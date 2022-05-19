// SPDX-License-Identifier: BUSL-1.1

package usereconomy

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/framey-io/go-tarantool"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
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

	err := s.createOrUpdateUserEconomy(u.User)
	if err != nil {
		return errors.Wrapf(err, "unable to call createOrUpdateUserEconomy")
	}

	return errors.Wrap(s.createEarnings(u.User), "unable to call createEarnings")
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
	space := s.totalUsersSpace()
	ix := s.totalUsersSpacePKIndex()
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
	nowT := uint64(time.Now().UTC().UnixNano())

	ue := &userEconomy{
		UserID:              u.ID,
		Username:            u.Username,
		ProfilePictureURL:   u.ProfilePictureURL,
		Balance:             0.0,
		StakingPercentage:   0.0,
		HashCode:            u.HashCode,
		LastMiningStartedAt: 0,
		StakingYears:        0,
		CreatedAt:           nowT,
		UpdatedAt:           nowT,
		BalanceUpdatedAt:    0,
	}

	return errors.Wrapf(s.db.InsertTyped(space, ue, &[]*userEconomy{}),
		"failed to insert user economy record for user.ID:%v", u.ID)
}

func (s *userEconomySource) updateUserEconomy(ue *userEconomy) error {
	nowT := uint64(time.Now().UTC().UnixNano())
	space := s.totalUsersSpace()
	index := s.totalUsersSpacePKIndex()
	key := tarantool.StringKey{S: ue.UserID}

	//nolint:gomnd // Those are not magic numbers, those are the indexes of the fields.
	ops := []tarantool.Op{
		{Op: "=", Field: 1, Arg: ue.Username},
		{Op: "=", Field: 2, Arg: ue.ProfilePictureURL},
		{Op: "=", Field: 9, Arg: nowT},
	}

	return errors.Wrapf(s.db.UpdateTyped(space, index, key, ops, &[]*userEconomy{}),
		"failed to update user economy record for user.ID:%v", ue.UserID)
}

func (s *userEconomySource) createEarnings(u *user) error {
	t1, err := s.findReferralOf(s.t1ReferralEarningsSpace(), u.ReferredBy)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return errors.Wrapf(s.initializeReferralEarnings(s.t1ReferralEarningsSpace(), u.ReferredBy, u.ID),
				"unable to create referral earnings record")
		}

		return errors.Wrapf(err, "unable to call getReferralEarnings")
	}

	var result error

	if err = s.initializeReferralEarnings(s.t1ReferralEarningsSpace(), u.ReferredBy, u.ID); err != nil {
		result = multierror.Append(result, errors.Wrapf(err,
			"unable to initialize T1 referral earnings for user.ID:%v and referral.ID:%v ", u.ReferredBy, u.ID))
	}
	if err = s.initializeReferralEarnings(s.t2ReferralEarningsSpace(), t1, u.ID); err != nil {
		result = multierror.Append(result, errors.Wrapf(err,
			"unable to initialize T2 referral earnings for user.ID:%v and referral.ID:%v ", t1, u.ID))
	}

	return errors.Wrapf(result, "unable to create earnings")
}

func (s *userEconomySource) findReferralOf(space string, referralID UserID) (UserID, error) {
	params := map[string]interface{}{
		"referralID": referralID,
	}

	sql := fmt.Sprintf(`
		SELECT user_id 
		FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" 
		WHERE referral_user_id = :referralID`, space)

	var res []*referredBy
	if err := s.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return "", errors.Wrapf(err, "failed to get %q record for referralID:%v", space, referralID)
	}

	if len(res) == 0 {
		return "", errors.Wrapf(storage.ErrNotFound,
			"unable to find %q record for referralID:%v", space, referralID)
	}

	return res[0].UserID, nil
}

func (s *userEconomySource) initializeReferralEarnings(space string, userID, referral UserID) error {
	nowT := uint64(time.Now().UTC().UnixNano())

	earning := &referralEarnings{
		UserID:         userID,
		ReferralUserID: referral,
		Earnings:       0.0,
		CreatedAt:      nowT,
		UpdatedAt:      nowT,
	}

	return errors.Wrapf(s.db.InsertTyped(space, earning, &[]*referralEarnings{}),
		"failed to create %s record for user.ID:%v and referral.ID:%v", space, userID, referral)
}

func (s *userEconomySource) t1ReferralEarningsSpace() string {
	return "t1_referral_earnings"
}

func (s *userEconomySource) t2ReferralEarningsSpace() string {
	return "t2_referral_earnings"
}

func (s *userEconomySource) userEconomySpace() string {
	return "user_economy"
}

func (s *userEconomySource) userEconomySpacePKIndex() string {
	return fmt.Sprintf("pk_unnamed_%s_1", strings.ToUpper(s.userEconomySpace()))
}

func (s *userEconomySource) totalUsersSpace() string {
	return "total_users"
}

func (s *userEconomySource) totalUsersSpacePKIndex() string {
	return fmt.Sprintf("pk_unnamed_%s_1", strings.ToUpper(s.totalUsersSpace()))
}
