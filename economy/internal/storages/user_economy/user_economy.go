// SPDX-License-Identifier: BUSL-1.1

package usereconomy

import (
	"context"
	"fmt"
	"math"
	"strings"

	"github.com/framey-io/go-tarantool"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
)

func New(db tarantool.Connector) messagebroker.Processor {
	return &userEconomySource{db: db}
}

func (s *userEconomySource) Process(ctx context.Context, m *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	u := new(users.UserSnapshot)
	if err := json.Unmarshal(m.Value, u); err != nil {
		return errors.Wrapf(err, "userEconomySource: cannot unmarshall %v into %#v", string(m.Value), u)
	}

	if u.User == nil && u.Before != nil {
		if err := s.deleteUserEconomy(u.Before); err != nil {
			return errors.Wrap(err, "unable to call deleteUserEconomy")
		}

		return errors.Wrap(s.updateTotalUsers(-1), "unable to call updateTotalUsers")
	}

	return errors.Wrap(s.initializeEconomy(u.User), "unable to call initializeTables")
}

func (s *userEconomySource) initializeEconomy(u *users.User) error {
	if err := s.createOrUpdateUserEconomy(u); err != nil {
		return errors.Wrap(err, "unable to call createOrUpdateUserEconomy")
	}
	if err := s.createUserStaking(u); err != nil {
		return errors.Wrap(err, "unable to call createUserStaking")
	}
	if err := s.createGeneralEarnings(u); err != nil {
		return errors.Wrap(err, "unable to call createGeneralEarnings")
	}

	return errors.Wrap(s.createReferralEarnings(u), "unable to call createReferralEarnings")
}

func (s *userEconomySource) deleteUserEconomy(u *users.User) error {
	params := map[string]interface{}{
		"userId": u.ID,
	}

	sql := fmt.Sprintf("DELETE FROM %[1]v WHERE user_id = :userId", s.userEconomySpace())

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, params)),
		"failed to delete user economy record for u.ID:%v", u.ID)
}

func (s *userEconomySource) updateTotalUsers(diff int) error {
	space := s.globalSpace()

	op := "+"
	if math.Signbit(float64(diff)) {
		op = "-"
	}

	incrementOps := []tarantool.Op{
		{Op: op, Field: 1, Arg: diff},
	}

	return errors.Wrapf(s.db.UpsertAsync(space, &totalUsers{Value: 1, Key: "TOTAL_USERS"}, incrementOps).GetTyped(&[]*totalUsers{}),
		"failed to update %v record the KEY = 'TOTAL_USERS'", space)
}

func (s *userEconomySource) createOrUpdateUserEconomy(u *users.User) error {
	ue, err := s.getUserEconomy(u.ID)
	if err != nil {
		tErr := new(tarantool.Error)
		if errors.Is(err, storage.ErrNotFound) || (errors.As(err, tErr) && tErr.Code == tarantool.ER_TUPLE_NOT_FOUND) {
			if err = s.createUserEconomy(u); err != nil {
				return errors.Wrapf(err, "unable to call createUserEconomy")
			}

			return errors.Wrapf(s.updateTotalUsers(1), "unable to call updateTotalUsers")
		}

		return errors.Wrapf(err, "unable to call getUserEconomy")
	}

	ue.ProfilePictureURL = u.ProfilePictureURL

	return errors.Wrapf(s.updateUserEconomy(ue), "unable to call updateUserEconomy")
}

func (s *userEconomySource) getUserEconomy(userID UserID) (*userEconomy, error) {
	space := s.userEconomySpace()
	index := s.userEconomySpacePKIndex()
	key := tarantool.StringKey{S: userID}

	var res userEconomy
	if err := s.db.GetTyped(space, index, key, &res); err != nil {
		return nil, errors.Wrapf(err, "unable to get %q record for userID:%v", space, userID)
	}
	if res.UserID == "" {
		return nil, errors.Wrapf(storage.ErrNotFound, "not found user_economy %v", userID)
	}

	return &res, nil
}

func (s *userEconomySource) createUserEconomy(u *users.User) error {
	space := s.userEconomySpace()
	nowT := time.Now()

	ue := &userEconomy{
		CreatedAt:         nowT,
		UpdatedAt:         nowT,
		UserID:            u.ID,
		Username:          u.Username,
		ProfilePictureURL: u.ProfilePictureURL,
		HashCode:          u.HashCode,
	}

	return errors.Wrapf(s.db.InsertTyped(space, ue, &[]*userEconomy{}),
		"failed to insert user economy record for user.ID:%v", u.ID)
}

func (s *userEconomySource) createUserStaking(u *users.User) error {
	space := s.stakingSpace()
	nowT := time.Now()

	ue := &staking{
		CreatedAt:  nowT,
		UpdatedAt:  nowT,
		UserID:     u.ID,
		Percentage: 0,
		Years:      0,
	}

	return errors.Wrapf(s.db.InsertTyped(space, ue, &[]*staking{}),
		"failed to insert user economy record for user.ID:%v", u.ID)
}

func (s *userEconomySource) updateUserEconomy(ue *userEconomy) error {
	nowT := time.Now()
	space := s.userEconomySpace()
	index := s.userEconomySpacePKIndex()
	key := tarantool.StringKey{S: ue.UserID}

	//nolint:gomnd // Those are not magic numbers, those are the indexes of the fields.
	ops := []tarantool.Op{
		{Op: "=", Field: 2, Arg: nowT},
		{Op: "=", Field: 4, Arg: ue.Username},
		{Op: "=", Field: 5, Arg: ue.ProfilePictureURL},
	}

	return errors.Wrapf(s.db.UpdateTyped(space, index, key, ops, &[]*userEconomy{}),
		"failed to update user economy record for user.ID:%v", ue.UserID)
}

func (s *userEconomySource) createGeneralEarnings(u *users.User) error {
	var errs []error
	types := []string{balanceTypeStandard, balanceTypeStaking, balanceTypeTotal}
	for _, level := range []uint8{tierLevel0, tierLevel1, tierLevel2} {
		standard := generateGeneralBalanceType(balanceTypeStandard, level)
		staking := generateGeneralBalanceType(balanceTypeStaking, level)
		types = append(types, []string{standard, staking}...)
	}
	for _, t := range types {
		if err := s.initializeEarnings(u.ID, t); err != nil {
			errs = append(errs, errors.Wrapf(err, "unable to initialize %v balance for userID:%v", t, u.ID))
		}
	}

	return errors.Wrapf(multiErr(errs), "unable to initialize user earnings")
}

func (s *userEconomySource) createReferralEarnings(u *users.User) error {
	var errs []error
	if err := s.initializeReferralEarnings(u.ID, u.ReferredBy, tierLevel0); err != nil {
		errs = append(errs, errors.Wrapf(err,
			"unable to initialize T0 referral earnings for userID:%v and type~[userID]:%v", u.ReferredBy, u.ID))
	}

	nextRID := u.ReferredBy
	for _, level := range []TierLevel{tierLevel1, tierLevel2} {
		balanceType := generateUserBalanceType(nextRID, balanceTypeStandard, tierLevel0)
		referralID, err := s.findReferralOf(balanceType)
		if err != nil {
			return errors.Wrapf(multiErr(errs), "unable to find referral for type:%v", balanceType)
		}
		if referralID == "" {
			break
		}
		if err := s.initializeReferralEarnings(u.ID, referralID, level); err != nil {
			errs = append(errs, errors.Wrapf(err,
				"unable to initialize T%v referral earnings for userID:%v and type~[userID]:%v", level, referralID, u.ID))
		}
		nextRID = referralID
	}

	return errors.Wrapf(multiErr(errs), "unable to create referral earnings")
}

func (s *userEconomySource) initializeReferralEarnings(userID, referredBy UserID, tierLevel TierLevel) error {
	types := []string{
		generateUserBalanceType(userID, balanceTypeStandard, tierLevel),
		generateUserBalanceType(userID, balanceTypeStaking, tierLevel),
	}
	for _, t := range types {
		if err := s.initializeEarnings(referredBy, t); err != nil {
			return errors.Wrapf(err, "unable to initialize earnings for userID:%[1]v and type:%[2]v", referredBy, t)
		}
	}

	return nil
}

func (s *userEconomySource) findReferralOf(balanceType BalanceType) (UserID, error) {
	space := s.balancesSpace()
	params := map[string]interface{}{
		"type": balanceType,
	}

	sql := fmt.Sprintf(`
		SELECT user_id 
			FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1"
			WHERE type = :type`, space)

	var res []*tier
	if err := s.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return "", errors.Wrapf(err, "failed to get %q record for type:%v", space, balanceType)
	}

	if len(res) == 0 {
		return "", errors.Wrapf(storage.ErrNotFound,
			"unable to find %q record for type:%v", space, balanceType)
	}

	return res[0].UserID, nil
}

func (s *userEconomySource) initializeEarnings(referredBy, balanceType BalanceType) error {
	space := s.balancesSpace()
	earning := &balances{
		UpdatedAt: time.Now(),
		Amount:    coin.NewAmountUint64(0),
		UserID:    referredBy,
		Type:      balanceType,
		AmountW0:  0,
		AmountW1:  0,
		AmountW2:  0,
		AmountW3:  0,
	}

	return errors.Wrapf(s.db.InsertTyped(space, earning, &[]*balances{}),
		"failed to create %s record for user.ID:%v", space, referredBy)
}

func generateUserBalanceType(userID UserID, balanceType BalanceType, tierLevel TierLevel) string {
	return fmt.Sprintf("%v~%v", generateGeneralBalanceType(balanceType, tierLevel), userID)
}

func generateGeneralBalanceType(balanceType BalanceType, tierLevel TierLevel) string {
	return fmt.Sprintf("t%v_referral_%v_earnings", tierLevel, balanceType)
}

func (s *userEconomySource) balancesSpace() string {
	return "BALANCES"
}

func (s *userEconomySource) userEconomySpace() string {
	return "USER_ECONOMY"
}

func (s *userEconomySource) globalSpace() string {
	return "GLOBAL"
}

func (s *userEconomySource) stakingSpace() string {
	return "STAKING"
}

func (s *userEconomySource) userEconomySpacePKIndex() string {
	return fmt.Sprintf("pk_unnamed_%s_1", strings.ToUpper(s.userEconomySpace()))
}

func multiErr(errs []error) error {
	if len(errs) > 0 {
		nonNilErrs := make([]error, 0, len(errs))
		for _, e := range errs {
			if e != nil {
				nonNilErrs = append(nonNilErrs, e)
			}
		}
		if len(nonNilErrs) > 0 {
			return multierror.Append(nil, nonNilErrs...)
		}
	}

	return nil
}
