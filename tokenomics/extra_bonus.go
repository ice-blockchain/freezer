// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"sync"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) ClaimExtraBonus(ctx context.Context, ebs *ExtraBonusSummary) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	now := time.Now()
	bonus, err := r.getAvailableExtraBonus(ctx, now, ebs.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonus for userID:%v", ebs.UserID)
	}
	sql := `UPDATE extra_bonus_processing_worker
		    SET extra_bonus = $6,
		        extra_bonus_started_at = $3,
		        extra_bonus_ended_at = $4
          	WHERE worker_index = $1
              AND user_id = $2
              AND $5 > coalesce(extra_bonus_started_at,'1999-01-08 04:05:06'::timestamp)`
	const argCount = 6
	args := append(make([]any, 0, argCount),
		r.workerIndex(ctx),
		ebs.UserID,
		*now.Time,
		now.Add(r.cfg.ExtraBonuses.Duration),
		now.Add(-r.cfg.ExtraBonuses.ClaimWindow),
		bonus.AvailableExtraBonus)
	affectedRows, err := storage.Exec(ctx, r.db, sql, args...)
	if err != nil {
		return errors.Wrapf(err, "failed to update extra_bonus_processing_worker to claim bonus for args:%#v", args) //nolint:asasalint // Intended.
	}
	if affectedRows == 0 {
		return ErrDuplicate
	}
	*ebs = *bonus

	return nil
}

//nolint:funlen,lll // .
func (r *repository) getAvailableExtraBonus(ctx context.Context, now *time.Time, userID string) (*ExtraBonusSummary, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := `SELECT bal_worker.last_mining_started_at,
				   bal_worker.last_mining_ended_at,
				   eb_worker.news_seen,
				   b.bonus AS flat_bonus,
				   (100 - (25 *  ((CASE WHEN ($3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $8::bigint) / $11)) < $9::bigint THEN 0 ELSE ($3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $8::bigint) / $11)) END)/$10::bigint))) AS bonus_percentage_remaining,
				   $12 < coalesce(eb_worker.extra_bonus_started_at,'1999-01-08 04:05:06'::timestamp) AS already_claimed
			FROM extra_bonus_start_date sd
				JOIN extra_bonus_processing_worker eb_worker
				  ON eb_worker.worker_index = $1
				 AND eb_worker.user_id = $2
				JOIN extra_bonuses b 
				  ON b.ix = ($3::bigint + (eb_worker.utc_offset * $4::bigint) - sd.value) / $5::bigint
				 AND $3::bigint + (eb_worker.utc_offset * $4::bigint) > sd.value
				 AND b.bonus > 0
				JOIN extra_bonuses_worker e
				  ON e.worker_index = $1
				 AND b.ix = e.extra_bonus_index
				 AND $3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $8::bigint) / $11) < $7::bigint
				 AND $3::bigint + (eb_worker.utc_offset * $4::bigint) - (sd.value + (e.extra_bonus_index * $5::bigint)) - $6::bigint - ((e.offset_value * $8::bigint) / $11) > 0
				JOIN balance_recalculation_worker bal_worker
				  ON bal_worker.worker_index = $1
				 AND bal_worker.user_id = eb_worker.user_id
			WHERE sd.key = 0`
	const networkLagDelta, argCount = 1.3, 12
	args := append(make([]any, 0, argCount),
		r.workerIndex(ctx),
		userID,
		now.UnixNano(),
		r.cfg.ExtraBonuses.UTCOffsetDuration,
		r.cfg.ExtraBonuses.Duration,
		r.cfg.ExtraBonuses.TimeToAvailabilityWindow,
		r.cfg.ExtraBonuses.ClaimWindow,
		r.cfg.ExtraBonuses.AvailabilityWindow,
		stdlibtime.Duration(float64(r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds())*networkLagDelta),
		r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow,
		r.cfg.WorkerCount,
		now.Add(-r.cfg.ExtraBonuses.ClaimWindow))
	res, err := storage.Get[struct {
		LastMiningStartedAt, LastMiningEndedAt        *time.Time
		NewsSeen, FlatBonus, BonusPercentageRemaining uint64
		AlreadyClaimed                                bool
	}](ctx, r.db, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select for available extra bonus for userID:%v", userID)
	}
	if res.AlreadyClaimed {
		return nil, ErrDuplicate
	}

	return &ExtraBonusSummary{
		AvailableExtraBonus: r.calculateExtraBonus(res.FlatBonus, res.BonusPercentageRemaining, res.NewsSeen, r.calculateMiningStreak(now, res.LastMiningStartedAt, res.LastMiningEndedAt)), //nolint:lll // .
	}, nil
}

func (r *repository) calculateExtraBonus(flatBonus, bonusPercentageRemaining, newsSeen, miningStreak uint64) (extraBonus uint64) {
	if flatBonus == 0 {
		return 0
	}
	if miningStreak >= uint64(len(r.cfg.ExtraBonuses.MiningStreakValues)) {
		extraBonus += r.cfg.ExtraBonuses.MiningStreakValues[len(r.cfg.ExtraBonuses.MiningStreakValues)-1]
	} else {
		extraBonus += r.cfg.ExtraBonuses.MiningStreakValues[miningStreak]
	}
	if newsSeenBonusValues := r.cfg.ExtraBonuses.NewsSeenValues; newsSeen >= uint64(len(newsSeenBonusValues)) {
		extraBonus += newsSeenBonusValues[len(newsSeenBonusValues)-1]
	} else {
		extraBonus += newsSeenBonusValues[newsSeen]
	}

	return ((extraBonus + flatBonus) * bonusPercentageRemaining) / percentage100
}

func (r *repository) initializeExtraBonusWorkers(ctx context.Context) {
	allWorkers := make(map[int16]map[uint64]uint64, r.cfg.WorkerCount)
	for extraBonusIndex := 0; extraBonusIndex < len(r.cfg.ExtraBonuses.FlatValues); extraBonusIndex++ {
		offsets := make([]uint64, r.cfg.WorkerCount, r.cfg.WorkerCount) //nolint:gosimple // Prefer to be more descriptive.
		for i := 0; i < len(offsets); i++ {
			offsets[i] = uint64(i)
		}
		rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(offsets), func(i, j int) { //nolint:gosec // Not a problem here.
			offsets[i], offsets[j] = offsets[j], offsets[i]
		})
		for workerIndex := int16(0); workerIndex < r.cfg.WorkerCount; workerIndex++ {
			if _, found := allWorkers[workerIndex]; !found {
				allWorkers[workerIndex] = make(map[uint64]uint64, len(r.cfg.ExtraBonuses.FlatValues))
			}
			allWorkers[workerIndex][uint64(extraBonusIndex)] = offsets[workerIndex]
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(int(r.cfg.WorkerCount))
	for key, val := range allWorkers {
		go func(workerIndex int16, extraBonusesWorkerValues map[uint64]uint64) {
			defer wg.Done()
			r.mustPopulateExtraBonusWorker(ctx, workerIndex, extraBonusesWorkerValues)
		}(key, val)
	}
	wg.Wait()
}

func (r *repository) mustPopulateExtraBonusWorker(ctx context.Context, workerIndex int16, extraBonusesWorkerValues map[uint64]uint64) {
	const argOffset = 2
	ix := uint64(0)
	values := make([]string, 0, r.cfg.WorkerCount)
	args := append(make([]any, 0), workerIndex)
	for extraBonusIndex, offset := range extraBonusesWorkerValues {
		values = append(values, fmt.Sprintf("($1,$%v,$%v)", ix+argOffset, ix+argOffset+1))
		args = append(args, extraBonusIndex, offset)
		ix += argOffset
	}
	reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
	defer cancel()
	sql := fmt.Sprintf(`INSERT INTO extra_bonuses_worker (worker_index,extra_bonus_index,offset_value) 
													     VALUES %[1]v
							   ON CONFLICT (worker_index, extra_bonus_index) DO NOTHING`, strings.Join(values, ","))
	_, err := storage.Exec(reqCtx, r.db, sql, args...)
	log.Panic(errors.Wrapf(err, "failed to initialize extra_bonuses_%[1]v", workerIndex))
}

func (s *deviceMetadataTableSource) Process(ctx context.Context, msg *messagebroker.Message) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(msg.Value) == 0 {
		return nil
	}
	var dm deviceMetadata
	if err := json.UnmarshalContext(ctx, msg.Value, &dm); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &dm)
	}
	if dm.UserID == "" || dm.TZ == "" || (dm.Before != nil && dm.Before.TZ == dm.TZ) {
		return nil
	}
	duration, err := stdlibtime.ParseDuration(strings.Replace(dm.TZ+"m", ":", "h", 1))
	if err != nil {
		return errors.Wrapf(err, "invalid timezone:%#v", &dm)
	}
	var (
		workerIndex int16
		hashCode    int64
	)
	if err = retry(ctx, func() error {
		workerIndex, hashCode, err = s.getWorker(ctx, dm.UserID)

		return errors.Wrapf(err, "failed to getWorker for userID:%v", dm.UserID)
	}); err != nil {
		return errors.Wrapf(err, "permanently failed to getWorker for userID:%v", dm.UserID)
	}
	sql := `INSERT INTO extra_bonus_processing_worker (worker_index, user_id, hash_code, utc_offset)
											   VALUES ($1		   , $2		, $3	   , $4)
			ON CONFLICT (worker_index, user_id) 
					DO UPDATE 
						  SET utc_offset = EXCLUDED.utc_offset
					WHERE extra_bonus_processing_worker.utc_offset != EXCLUDED.utc_offset`
	_, err = storage.Exec(ctx, s.db, sql, workerIndex, dm.UserID, hashCode, int16(duration/stdlibtime.Minute))

	return errors.Wrapf(err, "failed to update users' timezone for %#v", &dm)
}

func (s *viewedNewsSource) Process(ctx context.Context, msg *messagebroker.Message) (err error) { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(msg.Value) == 0 {
		return nil
	}
	var vn viewedNews
	if err = json.UnmarshalContext(ctx, msg.Value, &vn); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &vn)
	}
	if vn.UserID == "" {
		return nil
	}
	if _, err = storage.Exec(ctx, s.db, `INSERT INTO processed_seen_news (user_id,news_id) VALUES ($1, $2)`, vn.UserID, vn.NewsID); err != nil {
		return errors.Wrapf(err, "failed to insert PROCESSED_SEEN_NEWS:%#v", &vn)
	}
	var (
		workerIndex int16
		hashCode    int64
	)
	if err = retry(ctx, func() error {
		workerIndex, hashCode, err = s.getWorker(ctx, vn.UserID)

		return errors.Wrapf(err, "failed to getWorker for userID:%v", vn.UserID)
	}); err != nil {
		return errors.Wrapf(err, "permanently failed to getWorker for userID:%v", vn.UserID)
	}
	sql := `INSERT INTO extra_bonus_processing_worker (worker_index, user_id, hash_code, news_seen)
											   VALUES ($1		   , $2		, $3	   , 1)
			ON CONFLICT (worker_index, user_id) 
					DO UPDATE 
						  SET news_seen = extra_bonus_processing_worker.news_seen + 1`
	_, err = storage.Exec(ctx, s.db, sql, workerIndex, vn.UserID, hashCode)

	return errors.Wrapf(err, "failed to update users' newsSeen count for %#v", &vn)
}
