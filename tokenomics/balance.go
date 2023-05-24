// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"sort"
	"strings"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) GetBalanceSummary( //nolint:lll // .
	ctx context.Context, userID string,
) (*BalanceSummary, error) {
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	res, err := storage.Get[struct {
		BalanceSoloField
		BalanceT0Field
		BalanceT1Field
		BalanceT2Field
		PreStakingBonusField
		PreStakingAllocationField
	}](ctx, r.db, SerializedUsersKey(id))
	if err != nil || len(res) == 0 {
		if err == nil {
			err = errors.Wrapf(ErrRelationNotFound, "missing state for id:%v", id)
		}

		return nil, errors.Wrapf(err, "failed to get balanceSummary for id:%v", id)
	}
	t1Standard, t1PreStaking := ApplyPreStaking(res[0].BalanceT0+res[0].BalanceT1, res[0].PreStakingAllocation, res[0].PreStakingBonus)
	t2Standard, t2PreStaking := ApplyPreStaking(res[0].BalanceT2, res[0].PreStakingAllocation, res[0].PreStakingBonus)
	soloStandard, soloPreStaking := ApplyPreStaking(res[0].BalanceSolo, res[0].PreStakingAllocation, res[0].PreStakingBonus)

	return &BalanceSummary{
		Balances: Balances[string]{
			Total:                  fmt.Sprint(soloStandard + soloPreStaking + t1Standard + t1PreStaking + t2Standard + t2PreStaking),
			TotalNoPreStakingBonus: fmt.Sprint(res[0].BalanceSolo + res[0].BalanceT0 + res[0].BalanceT1 + res[0].BalanceT2),
			Standard:               fmt.Sprint(soloStandard + t1Standard + t2Standard),
			PreStaking:             fmt.Sprint(soloPreStaking + t1PreStaking + t2PreStaking),
			T1:                     fmt.Sprint(t1Standard + t1PreStaking),
			T2:                     fmt.Sprint(t2Standard + t2PreStaking),
			TotalReferrals:         fmt.Sprint(t1Standard + t1PreStaking + t2Standard + t2PreStaking),
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
	for ix := stdlibtime.Duration(0); ix < stdlibtime.Duration(cap(typeDetails)); ix++ {
		date := start.Add((ix + stdlibtime.Duration(mappedOffset-mappedLimit)) * factor * r.cfg.GlobalAggregationInterval.Child)
		typeDetails = append(typeDetails, fmt.Sprintf("/%v", date.Format(r.cfg.globalAggregationIntervalChildDateFormat())))
	}
	if true {
		return make([]*BalanceHistoryEntry, 0, 0), nil //nolint:gosimple // Nope.
	}
	adoptions, gErr := getAllAdoptions[float64](ctx, r.db)
	if gErr != nil {
		return nil, errors.Wrap(gErr, "failed to getAllAdoptions")
	}
	location := stdlibtime.FixedZone(utcOffset.String(), int(utcOffset.Seconds()))
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
		date = date.In(location)
		childDateFormat, parentDateFormat := date.Format(childDateLayout), date.Format(parentDateLayout)
		if _, found := filteredChildrenByParents[parentDateFormat]; !found {
			filteredChildrenByParents[parentDateFormat] = make(map[string]any, mappedLimit)
		}
		if _, found := filteredChildrenByParents[parentDateFormat][childDateFormat]; !found {
			filteredChildrenByParents[parentDateFormat][childDateFormat] = struct{}{}
		}
	}
	resp := make([]*BalanceHistoryEntry, 0, 1+1)
	for _, parent := range r.processBalanceHistory(nil, factor > 0, utcOffset, adoptions) {
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
	adoptions []*Adoption[float64],
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
					Balance: &BalanceHistoryBalanceDiff{},
				},
				children: make(map[string]*BalanceHistoryEntry, int(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)),
			}
		}
		if _, found := parents[parentFormat].children[childFormat]; !found {
			parents[parentFormat].children[childFormat] = &BalanceHistoryEntry{
				Time:    child.In(location),
				Balance: &BalanceHistoryBalanceDiff{},
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
		var baseMiningRate float64
		for _, childVal := range parentVal.children {
			childVal.Time = childVal.Time.Add(childMin30TzAdjustment).Add(childMin45TzAdjustment)
			baseMiningRate += childVal.calculateBalanceDiffBonus(r.cfg.GlobalAggregationInterval.Child, utcOffset, adoptions)
			parentVal.reduceBalance(childVal.Balance.Negative, childVal.Balance.amount)
			if r.cfg.GlobalAggregationInterval.Child == stdlibtime.Hour && childVal.Time.Minute() != 0 {
				childVal.Time = childVal.Time.Add(-stdlibtime.Duration(childVal.Time.Minute()) * stdlibtime.Minute)
			}
			childVal.Balance.Amount = fmt.Sprint(childVal.Balance.amount)
			parentVal.BalanceHistoryEntry.TimeSeries = append(parentVal.BalanceHistoryEntry.TimeSeries, childVal)
		}
		parentVal.setBalanceDiffBonus(baseMiningRate / (float64(len(parentVal.children))))
		parentVal.Balance.Amount = fmt.Sprint(parentVal.Balance.amount)
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

func (e *BalanceHistoryEntry) reduceBalance(negative bool, amount float64) { //nolint:revive // Not an issue here.
	if negative != e.Balance.Negative {
		if amount > e.Balance.amount { //nolint:gocritic // Nope.
			e.Balance.Negative = negative
			e.Balance.amount -= e.Balance.amount
		} else if amount < e.Balance.amount {
			e.Balance.amount -= amount
		} else {
			e.Balance.Negative = false
			e.Balance.amount = 0.0
		}
	} else {
		e.Balance.amount += amount
	}
}

func (e *BalanceHistoryEntry) calculateBalanceDiffBonus( //nolint:funlen // .
	delta, utcOffset stdlibtime.Duration, adoptions []*Adoption[float64],
) (baseMiningRate float64) {
	endDate := e.Time.Add(delta)
	calculateProportionalBaseMiningRate := func(currentBaseMiningRate float64, startDate stdlibtime.Time) float64 {
		return (currentBaseMiningRate * float64(endDate.Sub(startDate))) / float64(delta)
	}

	for ix := len(adoptions) - 1; ix >= 0; ix-- {
		if adoptions[ix].AchievedAt == nil {
			continue
		}
		achievedAt := adoptions[ix].AchievedAt.Add(utcOffset)
		currentBaseMiningRate := adoptions[ix].BaseMiningRate
		if achievedAt.Before(e.Time.Add(stdlibtime.Nanosecond)) {
			if baseMiningRate == 0 {
				baseMiningRate = currentBaseMiningRate
			} else {
				baseMiningRate += calculateProportionalBaseMiningRate(currentBaseMiningRate, e.Time)
			}

			break
		}
		if achievedAt.Before(endDate) && achievedAt.After(e.Time.Add(-stdlibtime.Nanosecond)) {
			baseMiningRate += calculateProportionalBaseMiningRate(currentBaseMiningRate, achievedAt)
			endDate = achievedAt
		}
	}
	e.setBalanceDiffBonus(baseMiningRate)

	return baseMiningRate
}

func (e *BalanceHistoryEntry) setBalanceDiffBonus(baseMiningRate float64) {
	if e.Balance.Negative { //nolint:gocritic // Wrong.
		e.Balance.Bonus = int64((-100 * (baseMiningRate + e.Balance.amount)) / baseMiningRate)
	} else if e.Balance.amount <= baseMiningRate {
		e.Balance.Bonus = int64((-100 * (baseMiningRate - e.Balance.amount)) / baseMiningRate)
	} else {
		e.Balance.Bonus = int64((100 * (e.Balance.amount - baseMiningRate)) / baseMiningRate)
	}
}

//nolint:funlen // .
func (s *completedTasksSource) Process(ctx context.Context, message *messagebroker.Message) (err error) {
	if ctx.Err() != nil || len(message.Value) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	const requiredCompletedTasks, adoptionMultiplicationFactor = 6, 150
	var val struct {
		UserID         string `json:"userId" example:"edfd8c02-75e0-4687-9ac2-1ce4723865c4"`
		CompletedTasks uint64 `json:"completedTasks,omitempty" example:"3"`
	}
	if err = json.UnmarshalContext(ctx, message.Value, &val); err != nil || val.UserID == "" || val.CompletedTasks != requiredCompletedTasks {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(message.Value), &val)
	}
	duplGuardKey := fmt.Sprintf("completed_tasks_ice_prize_dupl_guards:%v", val.UserID)
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
	adoption, err := GetCurrentAdoption(ctx, s.db)
	if err != nil {
		return errors.Wrap(err, "failed to getCurrentAdoption")
	}
	id, err := s.getOrInitInternalID(ctx, val.UserID)
	if err != nil {
		return errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", val.UserID)
	}
	prize := adoption.BaseMiningRate * adoptionMultiplicationFactor

	return errors.Wrapf(s.db.HIncrByFloat(ctx, SerializedUsersKey(id), "balance_solo_pending", prize).Err(),
		"failed to incr balance_solo_pending for userID:%v by %v", val.UserID, prize)
}

//nolint:gomnd // .
func ApplyPreStaking(amount float64, preStakingAllocation, preStakingBonus uint16) (float64, float64) {
	standardAmount := (amount * float64(100-preStakingAllocation)) / 100
	preStakingAmount := (amount * float64(100+preStakingBonus) * float64(preStakingAllocation)) / 10000

	return standardAmount, preStakingAmount
}
