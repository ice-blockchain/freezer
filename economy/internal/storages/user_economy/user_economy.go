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
	u := new(userSnapshot)
	if err := json.Unmarshal(m.Value, u); err != nil {
		return errors.Wrapf(err, "userEconomySource: cannot unmarshall %v into %#v", string(m.Value), u)
	}

	if !u.User.DeletedAt.IsZero() {
		if err := s.deleteUserEconomy(u.User); err != nil {
			return errors.Wrapf(err, "unable to call deleteUserEconomy")
		}

		return errors.Wrapf(s.updateTotalUsers(-1), "unable to call updateTotalUsers")
	}

	if err := s.createOrUpdateUserEconomy(u.User); err != nil {
		return errors.Wrapf(err, "unable to call createOrUpdateUserEconomy")
	}
	if err := s.createUserEarnings(u.User); err != nil {
		return errors.Wrap(err, "unable to call createUserEarnings")
	}

	return errors.Wrap(s.createReferralEarnings(u.User), "unable to call createReferralEarnings")
}

func (s *userEconomySource) deleteUserEconomy(u *user) error {
	params := map[string]interface{}{
		"userId": u.ID,
	}

	sql := fmt.Sprintf("DELETE FROM %[1]v WHERE user_id = :userId", s.userEconomySpace())

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, params)),
		"failed to delete user economy record for user.ID:%v", u.ID)
}

func (s *userEconomySource) updateTotalUsers(diff int) error {
	space := s.globalSpace()
	ix := s.globalSpacePKIndex()
	key := tarantool.StringKey{S: "TOTAL_USERS"}

	op := "+"
	if math.Signbit(float64(diff)) {
		op = "-"
	}

	incrementOps := []tarantool.Op{
		{Op: op, Field: 1, Arg: diff},
	}

	return errors.Wrapf(s.db.UpdateTyped(space, ix, key, incrementOps, &[]*totalUsers{}),
		"failed to update %v record the KEY = 'TOTAL_USERS'", space)
}

func (s *userEconomySource) createOrUpdateUserEconomy(u *user) error {
	ue, err := s.getUserEconomy(u.ID)
	if err != nil {
		tErr := new(tarantool.Error)
		if errors.As(err, tErr) && tErr.Code == tarantool.ER_TUPLE_NOT_FOUND {
			if err = s.updateTotalUsers(1); err != nil {
				return errors.Wrapf(err, "unable to call updateTotalUsers")
			}

			return errors.Wrapf(s.createUserEconomy(u), "unable to call createUserEconomy")
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

	var res *userEconomy
	if err := s.db.GetTyped(space, index, key, &res); err != nil {
		return nil, errors.Wrapf(err, "unable to get %q record for userID:%v", space, userID)
	}

	return res, nil
}

func (s *userEconomySource) createUserEconomy(u *user) error {
	space := s.userEconomySpace()
	nowT := time.Now()

	ue := &userEconomy{
		UserID:            u.ID,
		Username:          u.Username,
		ProfilePictureURL: u.ProfilePictureURL,
		HashCode:          u.HashCode,
		CreatedAt:         nowT,
		UpdatedAt:         nowT,
	}

	return errors.Wrapf(s.db.InsertTyped(space, ue, &[]*userEconomy{}),
		"failed to insert user economy record for user.ID:%v", u.ID)
}

func (s *userEconomySource) updateUserEconomy(ue *userEconomy) error {
	nowT := uint64(time.Now().UTC().UnixNano())
	space := s.userEconomySpace()
	index := s.userEconomySpacePKIndex()
	key := tarantool.StringKey{S: ue.UserID}

	//nolint:gomnd // Those are not magic numbers, those are the indexes of the fields.
	ops := []tarantool.Op{
		{Op: "=", Field: 2, Arg: ue.Username},
		{Op: "=", Field: 3, Arg: ue.ProfilePictureURL},
		{Op: "=", Field: 5, Arg: nowT},
	}

	return errors.Wrapf(s.db.UpdateTyped(space, index, key, ops, &[]*userEconomy{}),
		"failed to update user economy record for user.ID:%v", ue.UserID)
}

func (s *userEconomySource) createUserEarnings(u *user) error {
	var errs error
	if err := s.initializeEarnings(u.ID, balanceTypeStandard); err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "unable to initialize %v balance for userID:%v", balanceTypeStandard, u.ID))
	}
	if err := s.initializeEarnings(u.ID, balanceTypeStaking); err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "unable to initialize %v balance for userID:%v", balanceTypeStaking, u.ID))
	}
	if err := s.initializeEarnings(u.ID, balanceTypeTotal); err != nil {
		errs = multierror.Append(errs, errors.Wrapf(err, "unable to initialize %v balance for userID:%v", balanceTypeTotal, u.ID))
	}

	return errors.Wrapf(errs, "unable to initialize user earnings")
}

// TODO: check multierrors. Maybe it should be implemented by []error.
func (s *userEconomySource) createReferralEarnings(u *user) error {
	var errs error
	if err := s.initializeReferralEarningsByLevel(u.ID, u.ReferredBy, tierLevel0); err != nil {
		multierror.Append(errs, errors.Wrapf(err,
			"unable to initialize T0 referral earnings for userID:%v and type~[userID]:%v", u.ReferredBy, u.ID))
	}

	rID := u.ReferredBy
	for _, level := range []uint64{tierLevel1, tierLevel2} {
		balanceType := generateBalanceTypeWithUserID(rID, balanceTypeStandard, tierLevel0)
		rID, err := s.findReferralOf(balanceType)
		if err != nil {
			return errors.Wrapf(errs, "unable to find referral for type:%v", balanceType)
		}
		if err := s.initializeReferralEarningsByLevel(u.ID, rID, level); err != nil {
			multierror.Append(errs, errors.Wrapf(err,
				"unable to initialize T%v referral earnings for userID:%v and type~[userID]:%v", level, rID, u.ID))
		}
	}

	return errors.Wrapf(errs, "unable to create referral earnings")
}

func (s *userEconomySource) initializeReferralEarningsByLevel(userID, referredBy UserID, tierLevel uint64) error {
	types := []string{generateBalanceTypeWithUserID(userID, balanceTypeStandard, tierLevel),
		generateBalanceTypeWithUserID(userID, balanceTypeStaking, tierLevel),
		generateBalanceType(balanceTypeStandard, tierLevel),
		generateBalanceType(balanceTypeStaking, tierLevel)}

	for _, t := range types {
		if err := s.initializeEarnings(referredBy, t); err != nil {
			return errors.Wrapf(err,
				"unable to initialize earnings for userID:%[1]v and type:%[2]v", referredBy, t)
		}
	}

	return nil
}

func (s *userEconomySource) findReferralOf(balanceType string) (UserID, error) {
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
		return "", errors.Wrapf(err, "failed to get %q record for type:%v", balanceType)
	}

	if len(res) == 0 {
		return "", errors.Wrapf(storage.ErrNotFound,
			"unable to find %q record for type:%v", space, balanceType)
	}

	return res[0].UserID, nil
}

func (s *userEconomySource) initializeEarnings(referredBy, balanceType string) error {
	space := s.balancesSpace()
	nowT := time.Now()
	earning := &referralEarnings{
		UserID:    referredBy,
		Type:      balanceType,
		UpdatedAt: nowT,
	}

	return errors.Wrapf(s.db.InsertTyped(space, earning, &[]*referralEarnings{}),
		"failed to create %s record for user.ID:%v", space, referredBy)
}

func generateBalanceTypeWithUserID(userID UserID, balanceType string, tierLevel uint64) string {
	return fmt.Sprintf("%v~%v", generateBalanceType(balanceType, tierLevel), userID)
}

func generateBalanceType(balanceType string, tierLevel uint64) string {
	return fmt.Sprintf("t%v_referral_%v_earnings", tierLevel, balanceType)
}

func (s *userEconomySource) balancesSpace() string {
	return "balances"
}

func (s *userEconomySource) userEconomySpace() string {
	return "user_economy"
}

func (s *userEconomySource) userEconomySpacePKIndex() string {
	return fmt.Sprintf("pk_unnamed_%s_1", strings.ToUpper(s.userEconomySpace()))
}

func (s *userEconomySource) globalSpace() string {
	return "global"
}

func (s *userEconomySource) globalSpacePKIndex() string {
	return fmt.Sprintf("pk_unnamed_%s_1", strings.ToUpper(s.globalSpace()))
}
