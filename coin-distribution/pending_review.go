// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	appcfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func NewRepository(ctx context.Context, _ context.CancelFunc) Repository {
	var localCfg config
	appcfg.MustLoadFromKey(applicationYamlKey, &localCfg)
	if localCfg.AlertSlackWebhook == "" {
		log.Panic("`alert-slack-webhook` is missing")
	}
	if localCfg.Environment == "" {
		log.Panic("`environment` is missing")
	}
	if localCfg.ReviewURL == "" {
		log.Panic("`review-url` is missing")
	}

	return &repository{
		db:  storage.MustConnect(ctx, ddl, applicationYamlKey),
		cfg: &localCfg,
	}
}

func (r *repository) CheckHealth(ctx context.Context) error {
	return errors.Wrap(r.db.Ping(ctx), "[health-check] failed to ping DB for coindistribution.repository")
}

func (r *repository) Close() error {
	return errors.Wrap(r.db.Close(), "failed to close db")
}

//nolint:funlen // .
func (r *repository) GetCoinDistributionsForReview(ctx context.Context, arg *GetCoinDistributionsForReviewArg) (*CoinDistributionsForReview, error) { //nolint:lll // .
	conditions, whereArgs := arg.where()
	sql := fmt.Sprintf(`SELECT * 
						FROM coin_distributions_pending_review 
						WHERE 1=1 
						  AND %[1]v
						ORDER BY %[2]v 
						LIMIT $2 OFFSET $1`, strings.Join(append(conditions, "1=1"), " AND "), strings.Join(append(arg.orderBy(), "internal_id asc"), ", "))
	result, err := storage.ExecMany[struct {
		*PendingReview
		Day        stdlibtime.Time
		InternalID uint64
	}](ctx, r.db, sql, append([]any{arg.Cursor, arg.Limit}, whereArgs...)...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select coin_distributions_pending_review for %#v", arg)
	}
	distributions := make([]*PendingReview, len(result)) //nolint:makezero // .
	for i, d := range result {
		d.PendingReview.Ice = float64(d.PendingReview.IceInternal) / 100
		distributions[i] = d.PendingReview
	}
	conditions, whereArgs = arg.totalsWhere()
	sql = fmt.Sprintf(`SELECT count(1) AS rows,
							   coalesce(sum(ice),0) AS ice 
					   FROM coin_distributions_pending_review 
					   WHERE 1=1
						 AND %[1]v`, strings.Join(append(conditions, "1=1"), " AND "))
	total, err := storage.ExecOne[struct {
		Rows uint64
		Ice  uint64
	}](ctx, r.db, sql, whereArgs...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select coin_distributions_pending_review totals for %#v", arg)
	}
	nextCursor := uint64(0)
	if len(result) == int(arg.Limit) {
		nextCursor = arg.Cursor + arg.Limit
	}

	return &CoinDistributionsForReview{
		Distributions: distributions,
		Cursor:        nextCursor,
		TotalRows:     total.Rows,
		TotalIce:      float64(total.Ice) / 100,
	}, nil
}

func (a *GetCoinDistributionsForReviewArg) orderBy() []string {
	res := make([]string, 0, 4)

	if a.IceOrderBy != "" {
		res = append(res, fmt.Sprintf("ice %v", a.IceOrderBy))
	}
	if a.ReferredByUsernameOrderBy != "" {
		res = append(res, fmt.Sprintf("referred_by_username %v", a.ReferredByUsernameOrderBy))
	}
	if a.CreatedAtOrderBy != "" {
		res = append(res, fmt.Sprintf("created_at %v", a.CreatedAtOrderBy))
	}
	if a.UsernameOrderBy != "" {
		res = append(res, fmt.Sprintf("username %v", a.UsernameOrderBy))
	}

	return res
}

func (a *GetCoinDistributionsForReviewArg) where() ([]string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	i := 3
	if referredByUsernameKeyword := a.ReferredByUsernameKeyword; referredByUsernameKeyword != "" {
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "!", "!!")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "%", "!%")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "_", "!_")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "[", "![")
		referredByUsernameKeyword = referredByUsernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("referred_by_username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(referredByUsernameKeyword))
		i++
	}
	if usernameKeyword := a.UsernameKeyword; usernameKeyword != "" {
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "!", "!!")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "%", "!%")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "_", "!_")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "[", "![")
		usernameKeyword = usernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(usernameKeyword))
	}

	return conditions, args
}

func (a *GetCoinDistributionsForReviewArg) totalsWhere() ([]string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	i := 1
	if referredByUsernameKeyword := a.ReferredByUsernameKeyword; referredByUsernameKeyword != "" {
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "!", "!!")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "%", "!%")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "_", "!_")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "[", "![")
		referredByUsernameKeyword = referredByUsernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("referred_by_username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(referredByUsernameKeyword))
		i++
	}
	if usernameKeyword := a.UsernameKeyword; usernameKeyword != "" {
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "!", "!!")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "%", "!%")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "_", "!_")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "[", "![")
		usernameKeyword = usernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(usernameKeyword))
	}

	return conditions, args
}

//nolint:funlen // .
func (r *repository) ReviewCoinDistributions(ctx context.Context, reviewerUserID string, decision string) error {
	const sqlToCheckIfAnythingNeedsApproving = "SELECT true AS bogus WHERE exists (select 1 FROM coin_distributions_pending_review LIMIT 1)"
	switch strings.ToLower(decision) {
	case "approve":
		return storage.DoInTransaction(ctx, r.db, func(conn storage.QueryExecer) error {
			if _, err := storage.ExecOne[struct{ Bogus bool }](ctx, conn, sqlToCheckIfAnythingNeedsApproving); err != nil {
				if storage.IsErr(err, storage.ErrNotFound) {
					err = nil
				}

				return errors.Wrap(err, "failed to check if any rows in coin_distributions_pending_review exist")
			}
			if _, err := storage.Exec(ctx, conn, "call approve_coin_distributions($1,false,true)", reviewerUserID); err != nil {
				return errors.Wrap(err, "failed to call approve_coin_distributions")
			}

			return errors.Wrap(r.sendCurrentCoinDistributionsAvailableForReviewAreApprovedSlackMessage(ctx),
				"failed to sendCurrentCoinDistributionsAvailableForReviewAreApprovedSlackMessage")
		})
	case "approve-and-process-immediately":
		return storage.DoInTransaction(ctx, r.db, func(conn storage.QueryExecer) error {
			if _, err := storage.ExecOne[struct{ Bogus bool }](ctx, conn, sqlToCheckIfAnythingNeedsApproving); err != nil {
				if storage.IsErr(err, storage.ErrNotFound) {
					err = nil
				}

				return errors.Wrap(err, "failed to check if any rows in coin_distributions_pending_review exist")
			}
			if _, err := storage.Exec(ctx, conn, "call approve_coin_distributions($1,true,true)", reviewerUserID); err != nil {
				return errors.Wrap(err, "failed to call approve_coin_distributions")
			}

			return errors.Wrap(r.sendCurrentCoinDistributionsAvailableForReviewAreApprovedToBeProcessedImmediatelySlackMessage(ctx),
				"failed to sendCurrentCoinDistributionsAvailableForReviewAreApprovedToBeProcessedImmediatelySlackMessage")
		})
	case "deny":
		return storage.DoInTransaction(ctx, r.db, func(conn storage.QueryExecer) error {
			if _, err := storage.ExecOne[struct{ Bogus bool }](ctx, conn, sqlToCheckIfAnythingNeedsApproving); err != nil {
				if storage.IsErr(err, storage.ErrNotFound) {
					err = nil
				}

				return errors.Wrap(err, "failed to check if any rows in coin_distributions_pending_review exist")
			}
			if _, err := storage.Exec(ctx, conn, "call deny_coin_distributions($1,true)", reviewerUserID); err != nil {
				return errors.Wrap(err, "failed to call deny_coin_distributions")
			}

			return errors.Wrap(r.sendCurrentCoinDistributionsAvailableForReviewAreDeniedSlackMessage(ctx),
				"failed to sendCurrentCoinDistributionsAvailableForReviewAreDeniedSlackMessage")
		})
	default:
		log.Panic(fmt.Sprintf("unknown decision:`%v`", decision))
	}

	return ctx.Err()
}

func (r *repository) CollectCoinDistributionsForReview(ctx context.Context, records []*ByEarnerForReview) error {
	if len(records) == 0 {
		return nil
	}
	const columns = 9
	values := make([]string, 0, len(records))
	args := make([]any, 0, len(records)*columns)
	for ix, record := range records {
		values = append(values, generateValuesSQLParams(ix, columns))
		args = append(args,
			record.CreatedAt.Time,
			record.CreatedAt.Time,
			record.InternalID,
			int64(record.Balance*100),
			record.Username,
			record.ReferredByUsername,
			record.UserID,
			record.EarnerUserID,
			record.EthAddress)
	}
	sql := fmt.Sprintf(`INSERT INTO coin_distributions_by_earner(created_at,day,internal_id,balance,username,referred_by_username,user_id,earner_user_id,eth_address) 
																 VALUES %v
						ON CONFLICT (day, user_id, earner_user_id) DO UPDATE
							SET 
								created_at = EXCLUDED.created_at,
								balance = EXCLUDED.balance,
								username = EXCLUDED.username,
								referred_by_username = EXCLUDED.referred_by_username,
								eth_address = EXCLUDED.eth_address`, strings.Join(values, ",\n"))
	_, err := storage.Exec(ctx, r.db, sql, args...)

	return errors.Wrapf(err, "failed to insert into coin_distributions_by_earner [%v]", len(records))
}

func generateValuesSQLParams(index, columns int) string {
	params := make([]string, 0, columns)
	for ii := 1; ii <= columns; ii++ {
		params = append(params, fmt.Sprintf("$%v", index*columns+ii))
	}

	return fmt.Sprintf("(%v)", strings.Join(params, ","))
}

func (r *repository) NotifyCoinDistributionCollectionCycleEnded(ctx context.Context) error {
	sql := `INSERT INTO global(key,value) 
					   VALUES ('coin_collector_latest_collecting_date',$1),
							  ('new_coin_distributions_pending','true'),
							  ('coin_collector_forced_execution','false')
				ON CONFLICT (key) DO UPDATE
					SET value = EXCLUDED.value`
	_, err := storage.Exec(ctx, r.db, sql, time.Now().Format(stdlibtime.DateOnly))

	return errors.Wrap(err, "failed to update global.value for coin_collector_latest_collecting_date to now and mark new_coin_distributions_pending")
}

func (r *repository) GetCollectorSettings(ctx context.Context) (*CollectorSettings, error) {
	sql := `SELECT (SELECT value::timestamp FROM global WHERE key = $1) 				AS coin_collector_latest_collecting_date,
				   (SELECT value::timestamp FROM global WHERE key = $2) 				AS coin_collector_start_date,
				   (SELECT value::timestamp FROM global WHERE key = $3) 				AS coin_collector_end_date,
				   coalesce((SELECT value::bool FROM global WHERE key = $4),false) 		AS coin_collector_enabled,
				   coalesce((SELECT value::bool FROM global WHERE key = $5),false) 		AS coin_collector_forced_execution,
				   coalesce((SELECT value::int FROM global WHERE key = $6),0) 			AS coin_collector_min_mining_streaks_required,
				   coalesce((SELECT value::int FROM global WHERE key = $7),0) 			AS coin_collector_start_hour,
				   coalesce((SELECT value::int FROM global WHERE key = $8),0) 			AS coin_collector_min_balance_required,
				   coalesce((SELECT value FROM global WHERE key = $9),'') 				AS coin_collector_denied_countries`
	val, err := storage.ExecOne[struct {
		CoinCollectorLatestCollectingDate     *time.Time
		CoinCollectorStartDate                *time.Time
		CoinCollectorEndDate                  *time.Time
		CoinCollectorDeniedCountries          string
		CoinCollectorMinMiningStreaksRequired int
		CoinCollectorStartHour                int
		CoinCollectorMinBalanceRequired       int
		CoinCollectorEnabled                  bool
		CoinCollectorForcedExecution          bool
	}](ctx, r.db, sql,
		"coin_collector_latest_collecting_date",
		"coin_collector_start_date",
		"coin_collector_end_date",
		"coin_collector_enabled",
		"coin_collector_forced_execution",
		"coin_collector_min_mining_streaks_required",
		"coin_collector_start_hour",
		"coin_collector_min_balance_required",
		"coin_collector_denied_countries")
	if err != nil {
		return nil, errors.Wrap(err, "failed to select info about GetCollectorSettings")
	}
	countries := strings.Split(strings.ToLower(val.CoinCollectorDeniedCountries), ",")
	mappedCountries := make(map[string]struct{}, len(countries))
	for ix := range countries {
		mappedCountries[countries[ix]] = struct{}{}
	}

	return &CollectorSettings{
		DeniedCountries:          mappedCountries,
		LatestDate:               val.CoinCollectorLatestCollectingDate,
		StartDate:                val.CoinCollectorStartDate,
		EndDate:                  val.CoinCollectorEndDate,
		MinBalanceRequired:       float64(val.CoinCollectorMinBalanceRequired),
		StartHour:                val.CoinCollectorStartHour,
		MinMiningStreaksRequired: uint64(val.CoinCollectorMinMiningStreaksRequired),
		Enabled:                  val.CoinCollectorEnabled,
		ForcedExecution:          val.CoinCollectorForcedExecution,
	}, nil
}

func tryPrepareCoinDistributionsForReview(ctx context.Context, db *storage.DB) error {
	return storage.DoInTransaction(ctx, db, func(conn storage.QueryExecer) error {
		if _, err := storage.ExecOne[struct{ Bogus bool }](ctx, conn, "SELECT true AS bogus FROM global WHERE key = 'new_coin_distributions_pending' FOR UPDATE SKIP LOCKED"); err != nil {
			if storage.IsErr(err, storage.ErrNotFound) {
				err = nil
			}

			return errors.Wrap(err, "failed to check if we should start preparing new coin distributions for review")
		}

		if _, err := storage.Exec(ctx, conn, "call prepare_coin_distributions_for_review(true)"); err != nil {
			return errors.Wrap(err, "failed to call prepare_coin_distributions_for_review")
		}

		if rowsDeleted, err := storage.Exec(ctx, conn, "DELETE FROM global where key = 'new_coin_distributions_pending'"); err != nil || rowsDeleted != 1 {
			if err == nil {
				err = errors.Errorf("expected 1 rowsDeleted, actual: %v", rowsDeleted)
			}

			return errors.Wrap(err, "failed to del global.key='new_coin_distributions_pending'")
		}

		return errors.Wrap(sendNewCoinDistributionsAvailableForReviewSlackMessage(ctx), "failed to sendNewCoinDistributionsAvailableForReviewSlackMessage")
	})
}

func startPrepareCoinDistributionsForReviewMonitor(ctx context.Context, db *storage.DB) {
	ticker := stdlibtime.NewTicker(30 * stdlibtime.Second) //nolint:gomnd // .
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, 10*stdlibtime.Minute) //nolint:gomnd // .
			log.Error(errors.Wrap(tryPrepareCoinDistributionsForReview(reqCtx, db), "failed to tryPrepareCoinDistributionsForReview"))
			cancel()
		case <-ctx.Done():
			return
		}
	}
}
