// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"math"
	"slices"
	"sort"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	"github.com/ice-blockchain/freezer/model"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetBalanceSummary( //nolint:lll // .
	ctx context.Context, userID string,
) (*BalanceSummary, error) {
	id, err := GetOrInitInternalID(ctx, r.db, userID, r.cfg.WelcomeBonusV2Amount)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	res, err := storage.Get[struct {
		model.UserIDField
		model.LatestDeviceField
		model.BalanceSoloField
		model.BalanceT0Field
		model.BalanceT1Field
		model.BalanceT2Field
		model.BalanceSoloEthereumField
		model.BalanceT0EthereumField
		model.BalanceT1EthereumField
		model.BalanceT2EthereumField
		model.BalanceSoloEthereumMainnetRewardPoolContributionField
		model.BalanceT0EthereumMainnetRewardPoolContributionField
		model.BalanceT1EthereumMainnetRewardPoolContributionField
		model.BalanceT2EthereumMainnetRewardPoolContributionField
		model.PreStakingBonusField
		model.PreStakingAllocationField
	}](ctx, r.db, model.SerializedUsersKey(id))
	if err != nil || len(res) == 0 {
		if err == nil {
			err = errors.Wrapf(ErrRelationNotFound, "missing state for id:%v", id)
		}

		return nil, errors.Wrapf(err, "failed to get balanceSummary for id:%v", id)
	}
	if r.isAdvancedTeamDisabled(res[0].LatestDevice) {
		res[0].BalanceT2 = 0
	}
	t1Standard, t1PreStaking := ApplyPreStaking(res[0].BalanceT0+res[0].BalanceT1, res[0].PreStakingAllocation, res[0].PreStakingBonus)
	t2Standard, t2PreStaking := ApplyPreStaking(res[0].BalanceT2, res[0].PreStakingAllocation, res[0].PreStakingBonus)
	soloStandard, soloPreStaking := ApplyPreStaking(res[0].BalanceSolo, res[0].PreStakingAllocation, res[0].PreStakingBonus)

	return &BalanceSummary{
		Balances: Balances[string]{
			Total:                              fmt.Sprintf(floatToStringFormatter, soloStandard+soloPreStaking+t1Standard+t1PreStaking+t2Standard+t2PreStaking),
			TotalNoPreStakingBonus:             fmt.Sprintf(floatToStringFormatter, res[0].BalanceSolo+res[0].BalanceT0+res[0].BalanceT1+res[0].BalanceT2),
			Standard:                           fmt.Sprintf(floatToStringFormatter, soloStandard+t1Standard+t2Standard),
			PreStaking:                         fmt.Sprintf(floatToStringFormatter, soloPreStaking+t1PreStaking+t2PreStaking),
			T1:                                 fmt.Sprintf(floatToStringFormatter, t1Standard+t1PreStaking),
			T2:                                 fmt.Sprintf(floatToStringFormatter, t2Standard+t2PreStaking),
			TotalReferrals:                     fmt.Sprintf(floatToStringFormatter, t1Standard+t1PreStaking+t2Standard+t2PreStaking),
			TotalMiningBlockchain:              fmt.Sprintf(floatToStringFormatter, res[0].BalanceSoloEthereum+res[0].BalanceT0Ethereum+res[0].BalanceT1Ethereum+res[0].BalanceT2Ethereum),                                                                                                                     //nolint:lll // .
			TotalMainnetRewardPoolContribution: fmt.Sprintf(floatToStringFormatter, res[0].BalanceSoloEthereumMainnetRewardPoolContribution+res[0].BalanceT0EthereumMainnetRewardPoolContribution+res[0].BalanceT1EthereumMainnetRewardPoolContribution+res[0].BalanceT2EthereumMainnetRewardPoolContribution), //nolint:lll // .
		},
	}, nil
}

func (r *repository) GetBalanceHistory( //nolint:funlen,gocognit,revive,gocyclo,cyclop,revive // Better to be grouped together.
	ctx context.Context, userID string, start, end *time.Time, utcOffset stdlibtime.Duration, limit, offset uint64,
) ([]*BalanceHistoryEntry, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	start, end = time.New(start.UTC()), time.New(end.UTC())
	var factor stdlibtime.Duration
	if start.After(*end.Time) {
		factor = -1
	} else {
		factor = 1
	}
	dates, notBeforeTime, notAfterTime := r.calculateDates(limit, offset, start, end, factor)
	if len(dates) == 0 {
		return []*BalanceHistoryEntry{}, nil
	}
	id, gErr := GetOrInitInternalID(ctx, r.db, userID, r.cfg.WelcomeBonusV2Amount)
	if gErr != nil {
		return nil, errors.Wrapf(gErr, "failed to getOrInitInternalID for userID:%v", userID)
	}
	balanceHistory, gErr := r.dwh.SelectBalanceHistory(ctx, id, dates)
	if gErr != nil {
		return nil, errors.Wrapf(gErr, "failed to SelectBalanceHistory for id:%v,createdAts:%#v", id, dates)
	}

	return r.processBalanceHistory(balanceHistory, factor > 0, notBeforeTime, notAfterTime, utcOffset), nil
}

func (r *repository) calculateDates(limit, offset uint64, start, end *time.Time, factor stdlibtime.Duration) (dates []stdlibtime.Time, notBeforeTime, notAfterTime *time.Time) {
	const (
		hoursInADay = 24
	)
	var calculatedLimit, mappedOffset uint64
	var afterStartPadding, beforeStartPadding uint64
	if r.cfg.GlobalAggregationInterval.Child == stdlibtime.Minute {
		calculatedLimit = (limit / hoursInADay) * uint64(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)
		if factor > 0 {
			afterStartPadding = 60 - uint64(start.Add(stdlibtime.Duration(-calculatedLimit*uint64(stdlibtime.Minute))).Minute())
			beforeStartPadding = uint64(start.Add(stdlibtime.Duration(-calculatedLimit * uint64(stdlibtime.Minute))).Minute())
		} else {
			beforeStartPadding = uint64(start.Add(stdlibtime.Duration(-calculatedLimit*uint64(stdlibtime.Minute))).Minute()) + 1
		}
	} else {
		if offset > 0 { // Hack not to load extra records from FE.
			return nil, nil, nil
		}
	}
	mappedLimit := calculatedLimit + beforeStartPadding + afterStartPadding
	if r.cfg.GlobalAggregationInterval.Child == stdlibtime.Minute {
		mappedOffset = (offset / hoursInADay) * uint64(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)
	}
	dates = make([]stdlibtime.Time, 0, mappedLimit)

	if factor > 0 {
		if r.cfg.GlobalAggregationInterval.Child == stdlibtime.Minute {
			for ix := stdlibtime.Duration(mappedOffset); ix < stdlibtime.Duration(mappedLimit+mappedOffset); ix++ {
				dates = append(dates, start.Add(-stdlibtime.Duration(beforeStartPadding)*r.cfg.GlobalAggregationInterval.Child).Add(ix*factor*r.cfg.GlobalAggregationInterval.Child).Truncate(r.cfg.GlobalAggregationInterval.Child))
			}
			notBeforeTime = time.New(start.Add(stdlibtime.Duration(mappedOffset * uint64(stdlibtime.Minute))))
			notAfterTime = time.New(start.Add(stdlibtime.Duration((calculatedLimit + mappedOffset) * uint64(stdlibtime.Minute))))
		} else {
			notBeforeTime = start
			notAfterTime = end
			firstDayOfStartMonth := stdlibtime.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, stdlibtime.UTC)
			lastDayOfEndMonth := stdlibtime.Date(end.Year(), end.Month(), int(daysInMonth(end)), 0, 0, 0, 0, stdlibtime.UTC)
			for ix := 0; ix <= int(lastDayOfEndMonth.Sub(firstDayOfStartMonth).Hours()/hoursInADay); ix++ {
				dates = append(dates, firstDayOfStartMonth.Add(stdlibtime.Duration(ix)*hoursInADay*stdlibtime.Hour))
			}
		}
	} else {
		if r.cfg.GlobalAggregationInterval.Child == stdlibtime.Minute {
			for ix := stdlibtime.Duration(mappedOffset); ix < stdlibtime.Duration(mappedLimit+mappedOffset); ix++ {
				dates = append(dates, start.Add(ix*factor*r.cfg.GlobalAggregationInterval.Child).Truncate(r.cfg.GlobalAggregationInterval.Child))
			}
			notBeforeTime = time.New(start.Add(stdlibtime.Duration((-calculatedLimit - mappedOffset) * uint64(stdlibtime.Minute))))
			notAfterTime = time.New(start.Add(stdlibtime.Duration(-mappedOffset * uint64(stdlibtime.Minute))))
		} else {
			notBeforeTime = end
			notAfterTime = start
			firstDayOfEndMonth := stdlibtime.Date(end.Year(), end.Month(), 1, 0, 0, 0, 0, stdlibtime.UTC)
			lastDayOfStartMonth := stdlibtime.Date(start.Year(), start.Month(), int(daysInMonth(start)), 0, 0, 0, 0, stdlibtime.UTC)
			for ix := 0; ix <= int(lastDayOfStartMonth.Sub(firstDayOfEndMonth).Hours()/hoursInADay); ix++ {
				dates = append(dates, lastDayOfStartMonth.Add(-1*stdlibtime.Duration(ix)*hoursInADay*stdlibtime.Hour))
			}
		}
	}

	return
}

func (r *repository) processBalanceHistory(
	res []*dwh.BalanceHistory,
	startDateIsBeforeEndDate bool,
	notBeforeTime, notAfterTime *time.Time,
	utcOffset stdlibtime.Duration,
) []*BalanceHistoryEntry { //nolint:funlen,gocognit,revive // .
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
	childDateLayout := r.cfg.globalAggregationIntervalChildDateFormat()
	parentDateLayout := r.cfg.globalAggregationIntervalParentDateFormat()
	type parentType struct {
		*BalanceHistoryEntry
		children map[string]*BalanceHistoryEntry
	}
	parents := make(map[string]*parentType, 1+1)
	parentKeys := make([]string, 0, len(parents))
	for _, bal := range res {
		childFormat, parentFormat := bal.CreatedAt.Format(childDateLayout), bal.CreatedAt.Format(parentDateLayout)
		if _, found := parents[parentFormat]; !found {
			parent, pErr := stdlibtime.ParseInLocation(parentDateLayout, parentFormat, stdlibtime.UTC)
			log.Panic(pErr) //nolint:revive // Intended.
			parents[parentFormat] = &parentType{
				BalanceHistoryEntry: &BalanceHistoryEntry{
					Time:    parent.Add(-stdlibtime.Duration(utcOffset.Seconds()) * stdlibtime.Second).In(location),
					Balance: new(BalanceHistoryBalanceDiff),
				},
				children: make(map[string]*BalanceHistoryEntry, int(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)),
			}
			parentKeys = append(parentKeys, parentFormat)
		}
		if _, found := parents[parentFormat].children[childFormat]; !found {
			parents[parentFormat].children[childFormat] = &BalanceHistoryEntry{
				Time:    bal.CreatedAt.Add(-stdlibtime.Duration(utcOffset.Seconds()) * stdlibtime.Second).In(location),
				Balance: new(BalanceHistoryBalanceDiff),
			}
		}
		var total float64
		if bal.BalanceTotalSlashed >= 0 {
			total = bal.BalanceTotalMinted - bal.BalanceTotalSlashed
		} else {
			total = bal.BalanceTotalMinted - (bal.BalanceTotalSlashed * -1)
		}
		parents[parentFormat].children[childFormat].Balance.amount = total
		parents[parentFormat].children[childFormat].Balance.Negative = total < 0
		parents[parentFormat].children[childFormat].Balance.Amount = fmt.Sprintf(floatToStringFormatter, math.Abs(total))
	}
	sort.Strings(parentKeys)
	history := make([]*BalanceHistoryEntry, 0, len(parents))

	var prevParent *parentType
	var prevChild *BalanceHistoryEntry
	for _, pKey := range parentKeys {
		parents[pKey].BalanceHistoryEntry.TimeSeries = make([]*BalanceHistoryEntry, 0, len(parents[pKey].children))
		childrenKeys := make([]string, 0)
		for cKey := range parents[pKey].children {
			childrenKeys = append(childrenKeys, cKey)
		}
		sort.Strings(childrenKeys)
		for _, cKey := range childrenKeys {
			parents[pKey].children[cKey].TimeSeries = make([]*BalanceHistoryEntry, 0, 0)
			if prevChild == nil {
				parents[pKey].children[cKey].Balance.Bonus = 0
			} else {
				parents[pKey].children[cKey].setBalanceDiffBonus(prevChild.Balance.amount)
			}
			parents[pKey].Balance.amount += parents[pKey].children[cKey].Balance.amount
			if time.New(parents[pKey].children[cKey].Time).UnixNano() >= notBeforeTime.Add(-stdlibtime.Duration(utcOffset.Seconds())*stdlibtime.Second).In(location).UnixNano() &&
				time.New(parents[pKey].children[cKey].Time).UnixNano() <= notAfterTime.Add(-stdlibtime.Duration(utcOffset.Seconds())*stdlibtime.Second).In(location).UnixNano() {
				parents[pKey].BalanceHistoryEntry.TimeSeries = append(parents[pKey].BalanceHistoryEntry.TimeSeries, parents[pKey].children[cKey])
				prevChild = parents[pKey].children[cKey]
			}
			childrenKeys = append(childrenKeys, cKey)
		}
		if prevParent == nil {
			parents[pKey].Balance.Bonus = 0
		} else {
			parents[pKey].setBalanceDiffBonus(prevParent.Balance.amount)
		}
		parents[pKey].Balance.Negative = parents[pKey].Balance.amount < 0
		parents[pKey].Balance.Amount = fmt.Sprintf(floatToStringFormatter, math.Abs(parents[pKey].Balance.amount))
		if !startDateIsBeforeEndDate {
			sort.SliceStable(parents[pKey].BalanceHistoryEntry.TimeSeries, func(i, j int) bool {
				return parents[pKey].BalanceHistoryEntry.TimeSeries[i].Time.After(parents[pKey].BalanceHistoryEntry.TimeSeries[j].Time)
			})
		}
		if len(parents[pKey].BalanceHistoryEntry.TimeSeries) > 0 {
			history = append(history, parents[pKey].BalanceHistoryEntry)
		}
		prevParent = parents[pKey]
	}
	if !startDateIsBeforeEndDate {
		sort.SliceStable(history, func(i, j int) bool {
			return history[i].Time.After(history[j].Time)
		})
	}

	return history
}

func (e *BalanceHistoryEntry) setBalanceDiffBonus(from float64) {
	to := e.Balance.amount
	if from < 0 && to > 0 {
		e.Balance.Bonus = roundFloat64(((from - to) / from) * 100)
	} else if from != 0 {
		e.Balance.Bonus = roundFloat64(-1 * ((from - to) / from) * 100)
	}
}

//nolint:funlen // .
func (s *completedTasksSource) Process(ctx context.Context, message *messagebroker.Message) (err error) {
	if ctx.Err() != nil || len(message.Value) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if true {
		return nil
	}
	const requiredCompletedTasks, adoptionMultiplicationFactor = 6, 168
	var val struct {
		UserID         string  `json:"userId" example:"edfd8c02-75e0-4687-9ac2-1ce4723865c4"`
		Type           string  `json:"type" example:"claim_username"`
		CompletedTasks uint64  `json:"completedTasks,omitempty" example:"3"`
		Prize          float64 `json:"prize,omitempty" example:"100"`
	}
	if err = json.UnmarshalContext(ctx, message.Value, &val); err != nil || val.UserID == "" ||
		(!s.tasksV2Enabled(val.UserID) && val.CompletedTasks != requiredCompletedTasks) ||
		(s.tasksV2Enabled(val.UserID) && val.Type == "") {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(message.Value), &val)
	}
	var duplGuardKey string
	if s.tasksV2Enabled(val.UserID) {
		duplGuardKey = fmt.Sprintf("completed_tasks_ice_prize_dupl_guards:%v:%v", val.Type, val.UserID)
	} else {
		duplGuardKey = fmt.Sprintf("completed_tasks_ice_prize_dupl_guards:%v", val.UserID)
	}
	if set, dErr := s.db.SetNX(ctx, duplGuardKey, "", s.cfg.MiningSessionDuration.Min).Result(); dErr != nil || !set {
		if dErr == nil {
			dErr = ErrDuplicate
		}

		return errors.Wrapf(dErr, "SetNX failed for completed_tasks_ice_prize_dupl_guard, userID: %v", val.UserID)
	}
	defer func() {
		if err != nil {
			undoCtx, cancelUndo := context.WithTimeout(context.Background(), requestDeadline)
			defer cancelUndo()
			err = multierror.Append( //nolint:wrapcheck // .
				err,
				errors.Wrapf(s.db.Del(undoCtx, duplGuardKey).Err(), "failed to del completed_tasks_ice_prize_dupl_guard key"),
			).ErrorOrNil()
		}
	}()
	id, err := GetOrInitInternalID(ctx, s.db, val.UserID, s.cfg.WelcomeBonusV2Amount)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", val.UserID)
	}
	prize := val.Prize
	if !s.tasksV2Enabled(val.UserID) {
		res, err := storage.Get[struct{ model.CreatedAtField }](ctx, s.db, model.SerializedUsersKey(id))
		if err != nil || len(res) == 0 {
			if err == nil {
				err = errors.Wrapf(ErrRelationNotFound, "missing state for id:%v", id)
			}

			return errors.Wrapf(err, "failed to get GetAdoptionSummary for id:%v", id)
		}
		prize = s.cfg.BaseMiningRate(time.Now(), res[0].CreatedAt) * adoptionMultiplicationFactor
	}

	return errors.Wrapf(s.db.HIncrByFloat(ctx, model.SerializedUsersKey(id), "balance_solo_pending", prize).Err(),
		"failed to incr balance_solo_pending for userID:%v by %v", val.UserID, prize)
}

//nolint:gomnd // .
func ApplyPreStaking(amount, preStakingAllocation, preStakingBonus float64) (float64, float64) {
	standardAmount := (amount * (100 - preStakingAllocation)) / 100
	preStakingAmount := (amount * (100 + preStakingBonus) * preStakingAllocation) / 10000

	return standardAmount, preStakingAmount
}

func daysInMonth(t *time.Time) uint64 {
	y, m, _ := t.Date()

	return uint64(stdlibtime.Date(y, m+1, 0, 0, 0, 0, 0, stdlibtime.UTC).Day())
}

func (r *repository) tasksV2Enabled(userID string) bool {
	return r.cfg.TasksV2Enabled || (len(r.cfg.AdminUsers) > 0 && slices.Contains(r.cfg.AdminUsers, userID))
}
