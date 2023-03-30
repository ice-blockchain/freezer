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
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/go-tarantool-client"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) ClaimExtraBonus(ctx context.Context, ebs *ExtraBonusSummary) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	now := time.Now()
	bonus, err := r.getAvailableExtraBonus(ctx, now, ebs.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getAvailableExtraBonus for userID:%v", ebs.UserID)
	}
	params := make(map[string]any, 5) //nolint:gomnd // There's 5 keys there.
	params["user_id"] = ebs.UserID
	params["now_nanos"] = now
	params["extra_bonus"] = bonus.AvailableExtraBonus
	params["duration"] = r.cfg.ExtraBonuses.Duration
	params["claim_window"] = r.cfg.ExtraBonuses.ClaimWindow
	sql := fmt.Sprintf(`UPDATE extra_bonus_processing_worker_%[1]v
		  SET extra_bonus = :extra_bonus,
		      extra_bonus_started_at = :now_nanos,
		      extra_bonus_ended_at = :now_nanos + :duration
          WHERE user_id = :user_id
            AND :now_nanos - IFNULL(extra_bonus_started_at,0) > :claim_window`, r.workerIndex(ctx))
	if err = storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			err = ErrDuplicate
		}

		return errors.Wrapf(err, "failed to updated users to claim bonus for params:%#v", params)
	}
	*ebs = *bonus

	return nil
}

//nolint:funlen,lll // .
func (r *repository) getAvailableExtraBonus(ctx context.Context, now *time.Time, userID string) (*ExtraBonusSummary, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	const networkLagDelta = 1.3
	params := make(map[string]any, 10) //nolint:gomnd // .
	params["user_id"] = userID
	params["now_nanos"] = now
	params["duration"] = r.cfg.ExtraBonuses.Duration
	params["utc_offset_duration"] = r.cfg.ExtraBonuses.UTCOffsetDuration
	params["availability_window"] = r.cfg.ExtraBonuses.AvailabilityWindow
	params["delayed_claim_penalty_window"] = r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow
	params["first_delayed_claim_penalty_window"] = stdlibtime.Duration(float64(r.cfg.ExtraBonuses.DelayedClaimPenaltyWindow.Nanoseconds()) * networkLagDelta)
	params["time_to_availability_window"] = r.cfg.ExtraBonuses.TimeToAvailabilityWindow
	params["claim_window"] = r.cfg.ExtraBonuses.ClaimWindow
	params["worker_count"] = r.cfg.WorkerCount
	sql := fmt.Sprintf(`SELECT bal_worker.last_mining_started_at,
							   bal_worker.last_mining_ended_at,
							   eb_worker.news_seen,
							   b.bonus,
							   (100 - (25 *  ((CASE WHEN (:now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (e.extra_bonus_index * :duration)) - :time_to_availability_window  - ((e.offset * :availability_window) / :worker_count)) < :first_delayed_claim_penalty_window THEN 0 ELSE (:now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (e.extra_bonus_index * :duration)) - :time_to_availability_window  - ((e.offset * :availability_window) / :worker_count)) END)/:delayed_claim_penalty_window))) AS bonus_percentage_remaining,
							   :now_nanos - IFNULL(eb_worker.extra_bonus_started_at, 0) < :claim_window AS already_claimed
						FROM extra_bonus_start_date sd
							JOIN extra_bonus_processing_worker_%[1]v eb_worker
							  ON eb_worker.user_id = :user_id
						    JOIN extra_bonuses b 
							  ON b.ix = (:now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - sd.value) / :duration
							 AND :now_nanos + (eb_worker.utc_offset * :utc_offset_duration) > sd.value
							 AND b.bonus > 0
							JOIN extra_bonuses_%[1]v e
							  ON b.ix = e.extra_bonus_index
						     AND :now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (e.extra_bonus_index * :duration)) - :time_to_availability_window  - ((e.offset * :availability_window) / :worker_count) < :claim_window
						     AND :now_nanos + (eb_worker.utc_offset * :utc_offset_duration) - (sd.value + (e.extra_bonus_index * :duration)) - :time_to_availability_window  - ((e.offset * :availability_window) / :worker_count) > 0
							JOIN balance_recalculation_worker_%[1]v bal_worker
							  ON bal_worker.user_id = eb_worker.user_id
						WHERE sd.key = 0`, r.workerIndex(ctx))
	res := make([]*struct {
		_msgpack                                      struct{} `msgpack:",asArray"` //nolint:tagliatelle,revive,nosnakecase // To insert we need asArray
		LastMiningStartedAt, LastMiningEndedAt        *time.Time
		NewsSeen, FlatBonus, BonusPercentageRemaining uint64
		AlreadyClaimed                                bool
	}, 0, 1)
	if err := r.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to select for available extra bonus for userID:%v", userID)
	}
	if len(res) == 0 {
		return nil, ErrNotFound
	}
	if res[0].AlreadyClaimed {
		return nil, ErrDuplicate
	}

	return &ExtraBonusSummary{
		AvailableExtraBonus: r.calculateExtraBonus(res[0].FlatBonus, res[0].BonusPercentageRemaining, res[0].NewsSeen, r.calculateMiningStreak(now, res[0].LastMiningStartedAt, res[0].LastMiningEndedAt)), //nolint:lll // .
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

func (r *repository) initializeExtraBonusWorkers() {
	allWorkers := make(map[uint64]map[uint64]uint64, r.cfg.WorkerCount)
	for extraBonusIndex := 0; extraBonusIndex < len(r.cfg.ExtraBonuses.FlatValues); extraBonusIndex++ {
		offsets := make([]uint64, r.cfg.WorkerCount, r.cfg.WorkerCount) //nolint:gosimple // Prefer to be more descriptive.
		for i := 0; i < len(offsets); i++ {
			offsets[i] = uint64(i)
		}
		rand.New(rand.NewSource(time.Now().UnixNano())).Shuffle(len(offsets), func(i, j int) { //nolint:gosec // Not a problem here.
			offsets[i], offsets[j] = offsets[j], offsets[i]
		})
		for workerIndex := uint64(0); workerIndex < r.cfg.WorkerCount; workerIndex++ {
			if _, found := allWorkers[workerIndex]; !found {
				allWorkers[workerIndex] = make(map[uint64]uint64, len(r.cfg.ExtraBonuses.FlatValues))
			}
			allWorkers[workerIndex][uint64(extraBonusIndex)] = offsets[workerIndex]
		}
	}
	wg := new(sync.WaitGroup)
	wg.Add(int(r.cfg.WorkerCount))
	for key, val := range allWorkers {
		go func(workerIndex uint64, extraBonusesWorkerValues map[uint64]uint64) {
			defer wg.Done()
			r.mustPopulateExtraBonusWorker(workerIndex, extraBonusesWorkerValues)
		}(key, val)
	}
	wg.Wait()
}

func (r *repository) mustPopulateExtraBonusWorker(workerIndex uint64, extraBonusesWorkerValues map[uint64]uint64) {
	values := make([]string, 0, r.cfg.WorkerCount)
	for extraBonusIndex, offset := range extraBonusesWorkerValues {
		values = append(values, fmt.Sprintf("(%v,%v)", extraBonusIndex, offset))
	}
	sql := fmt.Sprintf("INSERT INTO extra_bonuses_%[1]v(extra_bonus_index,offset) VALUES %[2]v", workerIndex, strings.Join(values, ","))
	if err := storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, map[string]any{})); err != nil && !errors.Is(err, ErrDuplicate) {
		log.Panic(errors.Wrapf(err, "failed to initialize extra_bonuses_%[1]v", workerIndex))
	}
}

type (
	//nolint:govet // We can't update the order of fields in space.
	extraBonusProcessingWorker struct {
		_msgpack                               struct{} `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // .
		ExtraBonusStartedAt, ExtraBonusEndedAt *time.Time
		UserID                                 string
		UTCOffset                              int64
		NewsSeen, ExtraBonus                   uint64
		LastExtraBonusIndexNotified            *uint64
	}
)

func (s *deviceMetadataTableSource) Process(ctx context.Context, msg *messagebroker.Message) error {
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
	workerIndex, err := s.getWorkerIndex(ctx, dm.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getWorkerIndex for userID:%v", dm.UserID)
	}
	space := fmt.Sprintf("EXTRA_BONUS_PROCESSING_WORKER_%v", workerIndex)
	tuple := &extraBonusProcessingWorker{UserID: dm.UserID, UTCOffset: int64(duration / stdlibtime.Minute)}
	ops := append(make([]tarantool.Op, 0, 1), tarantool.Op{Op: "=", Field: 3, Arg: tuple.UTCOffset}) //nolint:gomnd // `utc_offset` column index.

	return errors.Wrapf(storage.CheckNoSQLDMLErr(s.db.UpsertTyped(space, tuple, ops, &[]*extraBonusProcessingWorker{})),
		"failed to update users' timezone for %#v", &dm)
}

func (s *viewedNewsSource) Process(ctx context.Context, msg *messagebroker.Message) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(msg.Value) == 0 {
		return nil
	}
	var vn viewedNews
	if err := json.UnmarshalContext(ctx, msg.Value, &vn); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &vn)
	}
	if vn.UserID == "" {
		return nil
	}
	if err := storage.CheckNoSQLDMLErr(s.db.InsertTyped("PROCESSED_SEEN_NEWS", &vn, &[]*viewedNews{})); err != nil {
		return errors.Wrapf(err, "failed to insert PROCESSED_SEEN_NEWS:%#v)", &vn)
	}
	workerIndex, err := s.getWorkerIndex(ctx, vn.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getWorkerIndex for userID:%v", vn.UserID)
	}
	space := fmt.Sprintf("EXTRA_BONUS_PROCESSING_WORKER_%v", workerIndex)
	tuple := &extraBonusProcessingWorker{UserID: vn.UserID, NewsSeen: 1}
	ops := append(make([]tarantool.Op, 0, 1), tarantool.Op{Op: "+", Field: 4, Arg: tuple.NewsSeen}) //nolint:gomnd // `news_seen` column index.
	if err = storage.CheckNoSQLDMLErr(s.db.UpsertTyped(space, tuple, ops, &[]*extraBonusProcessingWorker{})); err != nil {
		return multierror.Append( //nolint:wrapcheck // Not needed.
			errors.Wrapf(err, "failed to update users' newsSeen count for %#v", &vn),
			errors.Wrapf(storage.CheckNoSQLDMLErr(s.db.DeleteTyped("PROCESSED_SEEN_NEWS", "pk_unnamed_PROCESSED_SEEN_NEWS_1", []any{vn.UserID, vn.NewsID}, &[]*viewedNews{})), //nolint:lll // .
				"[rollback]failed to delete PROCESSED_SEEN_NEWS(%v,%v)", vn.UserID, vn.NewsID),
		).ErrorOrNil()
	}

	return nil
}
