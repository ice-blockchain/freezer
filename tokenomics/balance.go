// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"math"
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
	id, err := r.getOrInitInternalID(ctx, userID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to getOrInitInternalID for userID:%v", userID)
	}
	res, err := storage.Get[struct {
		model.UserIDField
		model.BalanceSoloField
		model.BalanceT0Field
		model.BalanceT1Field
		model.BalanceT2Field
		model.PreStakingBonusField
		model.PreStakingAllocationField
	}](ctx, r.db, model.SerializedUsersKey(id))
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
	ctx context.Context, userID string, start, end *time.Time, _ stdlibtime.Duration, limit, offset uint64,
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
	const (
		hoursInADay = 24
	)
	mappedLimit := (limit / hoursInADay) * uint64(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)
	mappedOffset := (offset / hoursInADay) * uint64(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)
	dates := make([]stdlibtime.Time, 0, mappedLimit)
	for ix := stdlibtime.Duration(mappedOffset); ix < stdlibtime.Duration(mappedLimit+mappedOffset); ix++ {
		dates = append(dates, start.Add(ix*factor*r.cfg.GlobalAggregationInterval.Child).Truncate(r.cfg.GlobalAggregationInterval.Child))
	}
	id, gErr := r.getOrInitInternalID(ctx, userID)
	if gErr != nil {
		return nil, errors.Wrapf(gErr, "failed to getOrInitInternalID for userID:%v", userID)
	}
	adoptions, gErr := getAllAdoptions[float64](ctx, r.db)
	if gErr != nil {
		return nil, errors.Wrap(gErr, "failed to getAllAdoptions")
	}
	balanceHistory, gErr := r.dwh.SelectBalanceHistory(ctx, id, dates)
	if gErr != nil {
		return nil, errors.Wrapf(gErr, "failed to SelectBalanceHistory for id:%v,createdAts:%#v", id, dates)
	}

	return r.processBalanceHistory(balanceHistory, factor > 0, adoptions), nil
}

func (r *repository) processBalanceHistory( //nolint:funlen,gocognit,revive // .
	res []*dwh.BalanceHistory,
	startDateIsBeforeEndDate bool,
	adoptions []*Adoption[float64],
) []*BalanceHistoryEntry {
	childDateLayout := r.cfg.globalAggregationIntervalChildDateFormat()
	parentDateLayout := r.cfg.globalAggregationIntervalParentDateFormat()
	parents := make(map[string]*struct {
		*BalanceHistoryEntry
		children map[string]*BalanceHistoryEntry
	}, 1+1)
	for _, bal := range res {
		childFormat, parentFormat := bal.CreatedAt.Format(childDateLayout), bal.CreatedAt.Format(parentDateLayout)
		if _, found := parents[parentFormat]; !found {
			parent, pErr := stdlibtime.ParseInLocation(parentDateLayout, parentFormat, stdlibtime.UTC)
			log.Panic(pErr) //nolint:revive // Intended.
			parents[parentFormat] = &struct {
				*BalanceHistoryEntry
				children map[string]*BalanceHistoryEntry
			}{
				BalanceHistoryEntry: &BalanceHistoryEntry{
					Time:    parent,
					Balance: new(BalanceHistoryBalanceDiff),
				},
				children: make(map[string]*BalanceHistoryEntry, int(r.cfg.GlobalAggregationInterval.Parent/r.cfg.GlobalAggregationInterval.Child)),
			}
		}
		if _, found := parents[parentFormat].children[childFormat]; !found {
			parents[parentFormat].children[childFormat] = &BalanceHistoryEntry{
				Time:    *bal.CreatedAt.Time,
				Balance: new(BalanceHistoryBalanceDiff),
			}
		}
		total := bal.BalanceTotalMinted - bal.BalanceTotalSlashed
		parents[parentFormat].children[childFormat].Balance.amount = total
		parents[parentFormat].children[childFormat].Balance.Negative = total < 0
		parents[parentFormat].children[childFormat].Balance.Amount = fmt.Sprint(math.Abs(total))
	}
	history := make([]*BalanceHistoryEntry, 0, len(parents))
	for _, parentVal := range parents {
		parentVal.BalanceHistoryEntry.TimeSeries = make([]*BalanceHistoryEntry, 0, len(parentVal.children))
		var baseMiningRate float64
		for _, childVal := range parentVal.children {
			baseMiningRate += childVal.calculateBalanceDiffBonus(r.cfg.GlobalAggregationInterval.Child, adoptions)
			parentVal.Balance.amount += childVal.Balance.amount
			parentVal.BalanceHistoryEntry.TimeSeries = append(parentVal.BalanceHistoryEntry.TimeSeries, childVal)
		}
		parentVal.setBalanceDiffBonus(baseMiningRate / (float64(len(parentVal.children))))
		parentVal.Balance.Negative = parentVal.Balance.amount < 0
		parentVal.Balance.Amount = fmt.Sprint(math.Abs(parentVal.Balance.amount))
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

func (e *BalanceHistoryEntry) calculateBalanceDiffBonus(
	delta stdlibtime.Duration, adoptions []*Adoption[float64],
) (baseMiningRate float64) {
	endDate := e.Time.Add(delta)
	calculateProportionalBaseMiningRate := func(currentBaseMiningRate float64, startDate stdlibtime.Time) float64 {
		return (currentBaseMiningRate * float64(endDate.Sub(startDate))) / float64(delta)
	}

	for ix := len(adoptions) - 1; ix >= 0; ix-- {
		if adoptions[ix].AchievedAt == nil {
			continue
		}
		achievedAt := *adoptions[ix].AchievedAt.Time
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

	return errors.Wrapf(s.db.HIncrByFloat(ctx, model.SerializedUsersKey(id), "balance_solo_pending", prize).Err(),
		"failed to incr balance_solo_pending for userID:%v by %v", val.UserID, prize)
}

//nolint:gomnd // .
func ApplyPreStaking(amount float64, preStakingAllocation, preStakingBonus uint16) (float64, float64) {
	standardAmount := (amount * float64(100-preStakingAllocation)) / 100
	preStakingAmount := (amount * float64(100+preStakingBonus) * float64(preStakingAllocation)) / 10000

	return standardAmount, preStakingAmount
}
