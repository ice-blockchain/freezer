// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"sort"
	"strings"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetBalanceSummary( //nolint:funlen,gocognit // Better to be grouped together.
	ctx context.Context, userID string,
) (*BalanceSummary, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`
SELECT b.*,
	   x.pre_staking_allocation,
	   st_b.bonus AS pre_staking_bonus,
       bal_worker.last_iteration_finished_at IS NOT NULL AND bal_worker.last_mining_ended_at IS NOT NULL
FROM (SELECT MAX(st.years) AS pre_staking_years,
		     MAX(st.allocation) AS pre_staking_allocation,
			 u.user_id,
			 u.referred_by
	  FROM users u
		 LEFT JOIN pre_stakings st
		        ON st.user_id = u.user_id
				AND st.worker_index = $1
      WHERE u.user_id = $2
	  GROUP BY u.user_id
	 ) x
   LEFT JOIN pre_staking_bonuses st_b
		  ON st_b.years = x.pre_staking_years
   JOIN balance_recalculation_worker bal_worker
		  ON bal_worker.user_id = x.user_id
		  AND bal_worker.worker_index = $3
   LEFT JOIN balances_worker b	
		  ON b.user_id = x.user_id
		 AND b.negative = FALSE
		 AND b.type = $4
		 AND b.type_detail IN ('','%[3]v_' || x.referred_by,'%[4]v','%[5]v')
		 AND b.worker_index = %[1]v`, r.workerIndex(ctx), totalNoPreStakingBonusBalanceType, t0BalanceTypeDetail, t1BalanceTypeDetail, t2BalanceTypeDetail) //nolint:lll // .
	type (
		B    = balance
		resp struct {
			*B
			PreStakingAllocation, PreStakingBonus uint64
			BalanceWorkerStarted                  bool
		}
	)
	res, err := storagev2.Select[resp](ctx, r.dbV2, sql, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select user's balances for user_id:%v", userID)
	}
	total, totalNoPreStakingBonus, t1, t2, standard, preStaking := coin.ZeroICEFlakes(), coin.ZeroICEFlakes(), coin.ZeroICEFlakes(), coin.ZeroICEFlakes(), coin.ZeroICEFlakes(), coin.ZeroICEFlakes() //nolint:lll // .
	for _, row := range res {
		if row.B == nil || row.B.Amount == nil {
			continue
		}
		standardAmount := row.Amount.
			MultiplyUint64(percentage100 - row.PreStakingAllocation).
			DivideUint64(percentage100)
		preStakingAmount := row.Amount.
			MultiplyUint64(row.PreStakingAllocation * (row.PreStakingBonus + percentage100)).
			DivideUint64(percentage100 * percentage100)
		switch row.TypeDetail {
		case t1BalanceTypeDetail:
			t1 = t1.Add(standardAmount.Add(preStakingAmount))
		case t2BalanceTypeDetail:
			t2 = standardAmount.Add(preStakingAmount)
		default:
			if strings.HasPrefix(row.TypeDetail, t0BalanceTypeDetail) {
				t1 = t1.Add(standardAmount.Add(preStakingAmount))
			}
		}
		standard = standard.Add(standardAmount)
		preStaking = preStaking.Add(preStakingAmount)
		total = total.Add(standardAmount.Add(preStakingAmount))
		totalNoPreStakingBonus = totalNoPreStakingBonus.Add(row.Amount)
	}
	if len(res) == 0 || !res[0].BalanceWorkerStarted { //nolint:revive // Wrong.
		standard = coin.NewAmountUint64(registrationICEFlakeBonusAmount)
		total = standard
		totalNoPreStakingBonus = total
	}

	return &BalanceSummary{
		Balances: Balances[coin.ICE]{
			Total:                  total.UnsafeICE(),
			TotalNoPreStakingBonus: totalNoPreStakingBonus.UnsafeICE(),
			Standard:               standard.UnsafeICE(),
			PreStaking:             preStaking.UnsafeICE(),
			T1:                     t1.UnsafeICE(),
			T2:                     t2.UnsafeICE(),
			TotalReferrals:         t1.Add(t2).UnsafeICE(),
		},
	}, nil
}

func (r *repository) GetBalanceHistory( //nolint:funlen,gocognit,revive,gocyclo,cyclop,revive // Better to be grouped together.
	ctx context.Context, userID string, start, end *time.Time, utcOffset stdlibtime.Duration, limit, offset uint64,
) ([]*BalanceHistoryEntry, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	var factor stdlibtime.Duration
	if start.After(*end.Time) {
		factor = -1
	} else {
		factor = 1
	}
	const (
		hoursInADay = 24
	)
	mappedLimit := (limit / hoursInADay) * uint64(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)
	mappedOffset := (offset / hoursInADay) * uint64(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)
	typeDetails := make([]string, 0, mappedLimit*2) //nolint:gomnd // Cuz we account for tz diff.
	args := make([]any, 0)
	args = append(args, userID, totalNoPreStakingBonusBalanceType)
	paramsCount := 1 + 1 + 1
	for ix := stdlibtime.Duration(0); ix < stdlibtime.Duration(cap(typeDetails)); ix++ {
		date := start.Add((ix + stdlibtime.Duration(mappedOffset-mappedLimit)) * factor * r.cfg.GlobalAggregationInterval.Child)
		args = append(args, fmt.Sprintf("/%v", date.Format(r.cfg.globalAggregationIntervalChildDateFormat())))
		typeDetails = append(typeDetails, fmt.Sprintf("$%d", paramsCount))
		paramsCount++
	}
	args = append(args, r.workerIndex(ctx))
	sql := fmt.Sprintf(`SELECT *
						FROM balances_worker
						WHERE user_id = $1
 					      AND (negative = TRUE OR negative = FALSE)
						  AND type = $2
						  AND type_detail in (%[1]v)
						  AND worker_index = $%[2]v`, strings.Join(typeDetails, ","), paramsCount)
	res, err := storagev2.Select[balance](ctx, r.dbV2, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select balance history for params:%#v", args...)
	}
	if len(res) == 0 {
		return make([]*BalanceHistoryEntry, 0, 0), nil //nolint:gosimple // Nope.
	}

	adoptions, gErr := getAllAdoptions[coin.ICEFlake](ctx, r.dbV2)
	if gErr != nil {
		return nil, errors.Wrap(gErr, "failed to getAllAdoptions")
	}

	preStakingSummaries, gErr := r.getAllPreStakingSummaries(ctx, userID)
	if gErr != nil {
		return nil, errors.Wrapf(gErr, "failed to getAllPreStakingSummaries for userID:%v", userID)
	}
	filteredChildrenByParents := make(map[string]map[string]any, 1+1)
	childDateLayout, parentDateLayout := r.cfg.globalAggregationIntervalChildDateFormat(), r.cfg.globalAggregationIntervalParentDateFormat()
	for ix := stdlibtime.Duration(mappedOffset); ix < stdlibtime.Duration(mappedLimit+mappedOffset); ix++ {
		date := start.Add((ix) * factor * r.cfg.GlobalAggregationInterval.Child)
		if factor == -1 && date.Before(*end.Time) {
			continue
		}
		if factor == 1 && date.After(*end.Time) {
			continue
		}
		date = date.Add(utcOffset)
		childDateFormat, parentDateFormat := date.Format(childDateLayout), date.Format(parentDateLayout)
		if _, found := filteredChildrenByParents[parentDateFormat]; !found {
			filteredChildrenByParents[parentDateFormat] = make(map[string]any, mappedLimit)
		}
		if _, found := filteredChildrenByParents[parentDateFormat][childDateFormat]; !found {
			filteredChildrenByParents[parentDateFormat][childDateFormat] = struct{}{}
		}
	}
	resp := make([]*BalanceHistoryEntry, 0, 1+1)
	for _, parent := range r.processBalanceHistory(res, factor > 0, utcOffset, adoptions, preStakingSummaries) {
		parentDateFormat := parent.Time.Format(parentDateLayout)
		if _, found := filteredChildrenByParents[parentDateFormat]; !found {
			continue
		}
		children := make([]*BalanceHistoryEntry, 0, len(parent.TimeSeries))
		for _, child := range parent.TimeSeries {
			if _, found := filteredChildrenByParents[parentDateFormat][child.Time.Format(childDateLayout)]; !found {
				continue
			}
			children = append(children, child)
		}
		if len(children) != 0 {
			parent.TimeSeries = children
			resp = append(resp, parent)
		}
	}

	return resp, nil
}

func (r *repository) processBalanceHistory( //nolint:funlen,gocognit,revive // .
	res []*balance,
	startDateIsBeforeEndDate bool,
	utcOffset stdlibtime.Duration,
	adoptions []*Adoption[coin.ICEFlake],
	preStakingSummaries []*PreStakingSummary,
) []*BalanceHistoryEntry {
	childDateLayout := r.cfg.globalAggregationIntervalChildDateFormat()
	parentDateLayout := r.cfg.globalAggregationIntervalParentDateFormat()
	parents := make(map[string]*struct {
		*BalanceHistoryEntry
		children map[string]*BalanceHistoryEntry
	}, 1+1)
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	for _, bal := range res {
		child, err := stdlibtime.Parse(childDateLayout, strings.Replace(bal.TypeDetail, "/", "", 1))
		log.Panic(err) //nolint:revive // Intended.
		child = child.In(location)
		childFormat, parentFormat := child.Format(childDateLayout), child.Format(parentDateLayout)
		if _, found := parents[parentFormat]; !found {
			parent, pErr := stdlibtime.Parse(parentDateLayout, parentFormat)
			log.Panic(pErr) //nolint:revive // Intended.
			parents[parentFormat] = &struct {
				*BalanceHistoryEntry
				children map[string]*BalanceHistoryEntry
			}{
				BalanceHistoryEntry: &BalanceHistoryEntry{
					Time:    parent.In(location),
					Balance: &BalanceHistoryBalanceDiff{amount: coin.ZeroICEFlakes()},
				},
				children: make(map[string]*BalanceHistoryEntry, int(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)),
			}
		}
		if _, found := parents[parentFormat].children[childFormat]; !found {
			parents[parentFormat].children[childFormat] = &BalanceHistoryEntry{
				Time:    child,
				Balance: &BalanceHistoryBalanceDiff{amount: coin.ZeroICEFlakes()},
			}
		}
		parents[parentFormat].children[childFormat].reduceBalance(bal.Negative, bal.Amount)
	}
	history := make([]*BalanceHistoryEntry, 0, len(parents))
	childMin30TzAdjustment, childMin45TzAdjustment := getTimezoneAdjustments(r.cfg.GlobalAggregationInterval.Child, utcOffset)
	parentMin30TzAdjustment, parentMin45TzAdjustment := getTimezoneAdjustments(r.cfg.GlobalAggregationInterval.Parent, utcOffset)
	for _, parentVal := range parents {
		parentVal.Time = parentVal.Time.Add(parentMin30TzAdjustment).Add(parentMin45TzAdjustment)
		parentVal.BalanceHistoryEntry.TimeSeries = make([]*BalanceHistoryEntry, 0, len(parentVal.children))
		var baseMiningRate *coin.ICEFlake
		for _, childVal := range parentVal.children {
			childVal.Time = childVal.Time.Add(childMin30TzAdjustment).Add(childMin45TzAdjustment)
			childVal.applyPreStaking(r.cfg.GlobalAggregationInterval.Child, utcOffset, preStakingSummaries)
			baseMiningRate = baseMiningRate.Add(childVal.calculateBalanceDiffBonus(r.cfg.GlobalAggregationInterval.Child, utcOffset, adoptions))
			parentVal.reduceBalance(childVal.Balance.Negative, childVal.Balance.amount)
			if r.cfg.GlobalAggregationInterval.Child == stdlibtime.Hour && childVal.Time.Minute() != 0 {
				childVal.Time = childVal.Time.Add(-stdlibtime.Duration(childVal.Time.Minute()) * stdlibtime.Minute)
			}
			childVal.Balance.Amount = childVal.Balance.amount.UnsafeICE()
			parentVal.BalanceHistoryEntry.TimeSeries = append(parentVal.BalanceHistoryEntry.TimeSeries, childVal)
		}
		parentVal.setBalanceDiffBonus(baseMiningRate.DivideUint64(uint64(len(parentVal.children))))
		parentVal.Balance.Amount = parentVal.Balance.amount.UnsafeICE()
		sort.SliceStable(parentVal.BalanceHistoryEntry.TimeSeries, func(i, j int) bool {
			if startDateIsBeforeEndDate {
				return parentVal.BalanceHistoryEntry.TimeSeries[i].Time.Before(parentVal.BalanceHistoryEntry.TimeSeries[j].Time)
			}

			return parentVal.BalanceHistoryEntry.TimeSeries[i].Time.After(parentVal.BalanceHistoryEntry.TimeSeries[j].Time)
		})
		history = append(history, parentVal.BalanceHistoryEntry)
	}
	sort.SliceStable(history, func(i, j int) bool {
		if startDateIsBeforeEndDate {
			return history[i].Time.Before(history[j].Time)
		}

		return history[i].Time.After(history[j].Time)
	})

	return history
}

func getTimezoneAdjustments(aggregationInterval, utcOffset stdlibtime.Duration) (min30Child, min45Child stdlibtime.Duration) {
	const halfHourTZFix = 30 * stdlibtime.Minute
	const min45TZFix = 45 * stdlibtime.Minute
	const min15TZFix = 15 * stdlibtime.Minute
	if aggregationInterval >= stdlibtime.Hour && utcOffset.Abs()%stdlibtime.Hour == halfHourTZFix {
		min30Child = -halfHourTZFix
	} else if aggregationInterval >= stdlibtime.Hour && utcOffset.Abs()%stdlibtime.Hour == min45TZFix {
		if utcOffset < 0 {
			min45Child = -min15TZFix
		} else {
			min45Child = -min45TZFix
		}
	}

	return
}

func (e *BalanceHistoryEntry) reduceBalance(negative bool, amount *coin.ICEFlake) { //nolint:revive // Not an issue here.
	if negative != e.Balance.Negative {
		if amount.GT(e.Balance.amount.Uint) { //nolint:gocritic // Nope.
			e.Balance.Negative = negative
			e.Balance.amount = amount.Subtract(e.Balance.amount)
		} else if amount.LT(e.Balance.amount.Uint) {
			e.Balance.amount = e.Balance.amount.Subtract(amount)
		} else {
			e.Balance.Negative = false
			e.Balance.amount = coin.ZeroICEFlakes()
		}
	} else {
		e.Balance.amount = e.Balance.amount.Add(amount)
	}
}

func (e *BalanceHistoryEntry) applyPreStaking( //nolint:funlen // .
	delta, utcOffset stdlibtime.Duration, preStakingSummaries []*PreStakingSummary,
) *BalanceHistoryEntry {
	if len(preStakingSummaries) == 0 {
		return e
	}
	var (
		resultingAmount = coin.ZeroICEFlakes()
		endDate         = e.Time.Add(delta)
	)
	applyProportionalPreStaking := func(ss *PreStakingSummary, startDate stdlibtime.Time) *coin.ICEFlake {
		return e.Balance.amount.
			MultiplyUint64(percentage100 - ss.Allocation).
			DivideUint64(percentage100).
			Add(e.Balance.amount.
				MultiplyUint64(ss.Allocation * (ss.Bonus + percentage100)).
				DivideUint64(percentage100 * percentage100)).
			MultiplyUint64(uint64(float64(endDate.Sub(startDate)) * coin.Denomination / float64(delta))).
			DivideUint64(coin.Denomination)
	}
	for ix := len(preStakingSummaries) - 1; ix >= 0; ix-- {
		preStakingCreatedAt := preStakingSummaries[ix].CreatedAt.Add(utcOffset)
		if preStakingCreatedAt.Before(e.Time.Add(stdlibtime.Nanosecond)) {
			resultingAmount = resultingAmount.Add(applyProportionalPreStaking(preStakingSummaries[ix], e.Time))

			break
		}
		if preStakingCreatedAt.Before(endDate) && preStakingCreatedAt.After(e.Time.Add(-stdlibtime.Nanosecond)) {
			resultingAmount = resultingAmount.Add(applyProportionalPreStaking(preStakingSummaries[ix], preStakingCreatedAt))
			endDate = preStakingCreatedAt
			if ix == 0 {
				resultingAmount = resultingAmount.Add(e.Balance.amount.
					MultiplyUint64(uint64(float64(endDate.Sub(e.Time)) * coin.Denomination / float64(delta))).
					DivideUint64(coin.Denomination))
			}
		}
	}
	if !resultingAmount.IsZero() {
		e.Balance.amount = resultingAmount
	}

	return e
}

func (e *BalanceHistoryEntry) calculateBalanceDiffBonus( //nolint:funlen // .
	delta, utcOffset stdlibtime.Duration, adoptions []*Adoption[coin.ICEFlake],
) (baseMiningRate *coin.ICEFlake) {
	endDate := e.Time.Add(delta)
	calculateProportionalBaseMiningRate := func(currentBaseMiningRate *coin.ICEFlake, startDate stdlibtime.Time) *coin.ICEFlake {
		return currentBaseMiningRate.
			MultiplyUint64(uint64(float64(endDate.Sub(startDate)) * coin.Denomination / float64(delta))).
			DivideUint64(coin.Denomination)
	}

	for ix := len(adoptions) - 1; ix >= 0; ix-- {
		if adoptions[ix].AchievedAt == nil {
			continue
		}
		achievedAt := adoptions[ix].AchievedAt.Add(utcOffset)
		currentBaseMiningRate := adoptions[ix].BaseMiningRate
		if achievedAt.Before(e.Time.Add(stdlibtime.Nanosecond)) {
			if baseMiningRate.IsZero() {
				baseMiningRate = currentBaseMiningRate
			} else {
				baseMiningRate = baseMiningRate.Add(calculateProportionalBaseMiningRate(currentBaseMiningRate, e.Time))
			}

			break
		}
		if achievedAt.Before(endDate) && achievedAt.After(e.Time.Add(-stdlibtime.Nanosecond)) {
			baseMiningRate = baseMiningRate.Add(calculateProportionalBaseMiningRate(currentBaseMiningRate, achievedAt))
			endDate = achievedAt
		}
	}
	e.setBalanceDiffBonus(baseMiningRate)

	return baseMiningRate
}

func (e *BalanceHistoryEntry) setBalanceDiffBonus(baseMiningRate *coin.ICEFlake) {
	if e.Balance.Negative { //nolint:gocritic // Wrong.
		e.Balance.Bonus = -1 * int64(baseMiningRate.
			Add(e.Balance.amount).
			MultiplyUint64(percentage100).
			Divide(baseMiningRate).Uint64())
	} else if e.Balance.amount.LTE(baseMiningRate.Uint) {
		e.Balance.Bonus = -1 * int64(baseMiningRate.
			Subtract(e.Balance.amount).
			MultiplyUint64(percentage100).
			Divide(baseMiningRate).Uint64())
	} else {
		e.Balance.Bonus = int64(e.Balance.amount.
			Subtract(baseMiningRate).
			MultiplyUint64(percentage100).
			Divide(baseMiningRate).Uint64())
	}
}

func (r *repository) insertOrReplaceBalances( //nolint:revive // Alot of SQL params and error handling. Control coupling is ok here.
	ctx context.Context, workerIndex int16, insert bool, updatedAt *time.Time, balances ...*balance,
) error {
	if ctx.Err() != nil || len(balances) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	values := make([]string, 0, len(balances))
	for _, bal := range balances {
		if bal.Amount.IsNil() {
			bal.Amount = coin.ZeroICEFlakes()
		}
		amount, err := bal.Amount.Uint.Marshal()
		log.Panic(err) //nolint:revive // Intended.
		values = append(values, fmt.Sprintf(`(%[1]v,'%[2]v',%[3]v,'%[4]v',%[5]v,'%[6]v','%[7]v','%[8]v')`,
			updatedAt.UnixNano(), string(amount), bal.UserID, bal.TypeDetail, bal.Type, bal.Negative, bal.HashCode, workerIndex))
	}

	sql := fmt.Sprintf(`INSERT INTO balances_worker (updated_at,amount,user_id,type_detail,type,negative,hash_code,worker_index)
			VALUES %v
			ON CONFLICT (user_id,type_detail,type,negative,worker_index) DO UPDATE SET
				updated_at=excluded.updated_at,
				amount=coalesce(excluded.amount, '0')`, strings.Join(values, ","))
	affectedRows, err := storagev2.Exec(ctx, r.dbV2, sql)
	if err != nil && affectedRows == 0 {
		return errors.Wrapf(err, "failed insert/update at %v balances:%#v for workerIndex:%v", updatedAt, balances, workerIndex)
	}

	return nil
}

func (r *repository) deleteBalances(ctx context.Context, workerIndex int16, balances ...*balance) error {
	if ctx.Err() != nil || len(balances) == 0 {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	values := make([]string, 0, len(balances))
	for _, bal := range balances {
		values = append(values, fmt.Sprintf(`(user_id = '%[1]v' AND negative = %[2]v AND type = %[3]v AND type_detail = '%[4]v')`,
			bal.UserID, bal.Negative, bal.Type, bal.TypeDetail))
	}
	sql := fmt.Sprintf(`DELETE FROM balances WHERE %v`, strings.Join(values, " OR "))
	affectedRows, err := storagev2.Exec(ctx, r.dbV2, sql)
	if err != nil && affectedRows == 0 {
		return errors.Wrapf(err, "failed to DELETE from balances for values:%#v", values)
	}

	return nil
}

func (r *repository) sendAddBalanceCommandMessage(ctx context.Context, cmd *AddBalanceCommand) error {
	valueBytes, err := json.MarshalContext(ctx, cmd)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal %#v", cmd)
	}
	msg := &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     cmd.UserID,
		Topic:   r.cfg.MessageBroker.Topics[5].Name,
		Value:   valueBytes,
	}
	responder := make(chan error, 1)
	defer close(responder)
	r.mb.SendMessage(ctx, msg, responder)

	return errors.Wrapf(<-responder, "failed to send `%v` message to broker", msg.Topic)
}

func (s *addBalanceCommandsSource) Process(ctx context.Context, message *messagebroker.Message) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(message.Value) == 0 {
		return nil
	}
	var val AddBalanceCommand
	if err := json.UnmarshalContext(ctx, message.Value, &val); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(message.Value), &val)
	}
	if val.UserID == "" {
		return nil
	}
	bal, err := s.balance(ctx, &val)
	if err != nil {
		return errors.Wrapf(err, "failed to build balance from %#v", val)
	}
	type processedAddBalanceCommand struct {
		_msgpack struct{} `msgpack:",asArray"` //nolint:unused,tagliatelle,revive,nosnakecase // To insert we need asArray
		UserID   string
		Key      string
	}
	sql := `INSERT INTO processed_add_balance_commands(user_id, key) VALUES($1, $2)
						ON CONFLICT (user_id, key) DO NOTHING`
	affectedRows, err := storagev2.Exec(ctx, s.dbV2, sql, val.UserID, val.EventID)
	if err != nil || affectedRows == 0 {
		return errors.Wrapf(err, "failed to insert PROCESSED_ADD_BALANCE_COMMANDS for userID:%v, key: %v)", val.UserID, val.EventID)
	}
	var (
		workerIndex int16
		hashCode    int64
	)
	if err = retry(ctx, func() error {
		workerIndex, hashCode, err = s.getWorker(ctx, val.UserID)

		return errors.Wrapf(err, "failed to getWorker for userID:%v", val.UserID)
	}); err != nil {
		return errors.Wrapf(err, "permanently failed to getWorker for userID:%v", val.UserID)
	}
	bal.HashCode, bal.WorkerIndex = hashCode, workerIndex
	err = errors.Wrapf(retry(ctx, func() error {
		if err = s.insertOrReplaceBalances(ctx, workerIndex, true, time.New(message.Timestamp), bal); err != nil {
			if errors.Is(err, storage.ErrRelationNotFound) {
				return err
			}

			return errors.Wrapf(backoff.Permanent(err), "failed to insertBalance:%#v", bal)
		}

		return nil
	}), "permanently failed to insertBalance:%#v", bal)
	if err != nil {
		sql := `DELETE FROM processed_add_balance_commands WHERE user_id = $1 AND key = $2`
		affectedRows, err := storagev2.Exec(ctx, s.dbV2, sql, val.UserID, val.EventID)
		if err != nil || affectedRows == 0 {
			return errors.Wrapf(err, "failed to delete PROCESSED_ADD_BALANCE_COMMAND(%v,%v)", val.UserID, val.EventID)
		}
	}

	return err
}

func (s *addBalanceCommandsSource) balance(ctx context.Context, cmd *AddBalanceCommand) (*balance, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	bal := &balance{
		UserID: cmd.UserID,
		Type:   pendingXBalanceType,
	}
	if cmd.Negative != nil && *cmd.Negative {
		bal.Negative = *cmd.Negative
	}
	if !cmd.T1.IsNil() {
		bal.Amount = cmd.T1
		bal.TypeDetail = t1BalanceTypeDetail
	}
	if !cmd.T2.IsNil() {
		bal.Amount = cmd.T2
		bal.TypeDetail = t2BalanceTypeDetail
	}
	if !cmd.Total.IsNil() {
		bal.Amount = cmd.Total
	}
	if !cmd.BaseFactor.IsNil() {
		adoption, err := s.getCurrentAdoption(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "failed to getCurrentAdoption")
		}
		bal.Amount = adoption.BaseMiningRate.Multiply(cmd.BaseFactor)
	}

	return bal, nil
}

func (b *balance) add(amount *coin.ICEFlake) {
	b.Amount = b.Amount.Add(amount)
}

func (b *balance) subtract(amount *coin.ICEFlake) {
	b.Amount = b.Amount.Subtract(amount)
}
