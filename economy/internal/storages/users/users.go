// SPDX-License-Identifier: BUSL-1.1

package users

import (
	"context"
	"fmt"
	"time"

	"github.com/framey-io/go-tarantool"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
)

func New(db tarantool.Connector) messagebroker.Processor {
	return &usersSource{db: db}
}

func (s *usersSource) Process(ctx context.Context, m *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	u := new(user)
	if err := json.Unmarshal(m.Value, u); err != nil {
		return errors.Wrapf(err, "usersSource: cannot unmarshall %v into %#v", string(m.Value), u)
	}

	if !u.DeletedAt.IsZero() {
		if err := s.deleteUserEconomy(u); err != nil {
			return errors.Wrapf(err, "unable to call deleteUserEconomy")
		}

		return errors.Wrapf(s.updateTotalUsers(-1), "unable to call updateTotalUsers")
	}

	err := s.createOrUpdateUserEconomy(u)
	if err != nil {
		return errors.Wrapf(err, "unable to call createOrUpdateUserEconomy")
	}

	return errors.Wrap(s.createEarnings(u), "unable to call createEarnings")
}

func (s *usersSource) updateTotalUsers(diff int) error {
	sql := fmt.Sprintf("UPDATE %[1]v SET"+
		"VALUE = VALUE %+d"+
		"WHERE KEY = 'TOTAL_USERS'", totalUsersSpace(), diff)

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, nil)),
		"failed to update %q record for the KEY = 'TOTAL_USERS'", totalUsersSpace())
}

func t1ReferralEarningsSpace() string {
	return "t1_referral_earnings"
}

func t2ReferralEarningsSpace() string {
	return "t2_referral_earnings"
}

func (s *usersSource) createEarnings(u *user) error {
	t1, err := s.findInT1(u.ReferredBy)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return errors.Wrapf(s.createReferral(t1ReferralEarningsSpace(), u.ReferredBy, u.ID),
				"unable to create referral earnings record")
		}

		return errors.Wrapf(err, "unable to call findInT1")
	}
	err1 := s.createReferral(t1ReferralEarningsSpace(), u.ReferredBy, u.ID)
	err2 := s.createReferral(t2ReferralEarningsSpace(), t1, u.ID)
	errs := make([]error, 0, 1+1)
	if err1 != nil {
		errs = append(errs, err1)
	}
	if err2 != nil {
		errs = append(errs, err2)
	}
	if len(errs) > 1 {
		return multierror.Append(nil, errs...)
	} else if len(errs) == 1 {
		return errors.Wrapf(errs[0], "failed to call createReferral")
	}

	return nil
}

func (s *usersSource) findInT1(referralID UserID) (UserID, error) {
	params := map[string]interface{}{
		"referralID": referralID,
	}

	sql := fmt.Sprintf(`
		SELECT user_id 
		FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" 
		WHERE referral_user_id = :referralID`, t1ReferralEarningsSpace())

	var res []*referredBy
	if err := s.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return "", errors.Wrapf(err, "failed to get %q record for referralID:%v", t1ReferralEarningsSpace(), referralID)
	}

	if len(res) == 0 {
		return "", errors.Wrapf(storage.ErrNotFound,
			"unable to find %q record for referralID:%v", t1ReferralEarningsSpace(), referralID)
	}

	return res[0].UserID, nil
}

func (s *usersSource) createReferral(space string, userID, referral UserID) error {
	nowT := uint64(time.Now().UTC().UnixNano())
	params := map[string]interface{}{
		"userID":     userID,
		"referralId": referral,
		"earnings":   0,
		"createdAt":  nowT,
		"updatedAt":  nowT,
	}

	sql := fmt.Sprintf("INSERT INTO %[1]v "+
		"(USER_ID, REFERRAL_USER_ID, EARININGS, CREATED_AT, UPDATED_AT) "+
		"VALUES "+
		"(:userID, :referralId, :earnings, :createdAt, :updatedAt)", space)

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, params)),
		"failed to create %s record for user.ID:%v and referral.ID:%v", space, userID, referral)
}

func (s *usersSource) deleteUserEconomy(u *user) error {
	params := map[string]interface{}{
		"userId": u.ID,
	}

	sql := fmt.Sprintf("DELETE FROM %[1]v WHERE user_id = :userId", userEconomySpace())

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, params)),
		"failed to delete user economy record for user.ID:%v", u.ID)
}

func (s *usersSource) createOrUpdateUserEconomy(u *user) error {
	ue, err := s.findUserEconomy(u.ID)
	switch {
	case errors.Is(err, storage.ErrNotFound):
		if err = s.updateTotalUsers(1); err != nil {
			return errors.Wrapf(err, "unable to call updateTotalUsers")
		}

		return errors.Wrapf(s.createUserEconomy(u), "unable to call createUserEconomy")
	case err != nil:
		return errors.Wrapf(err, "unable to call findUserEconomy")
	}

	ue.ProfilePictureURL = u.ProfilePictureURL

	return errors.Wrapf(s.updateUserEconomy(ue), "unable to call updateUserEconomy")
}

func (s *usersSource) findUserEconomy(userID UserID) (*userEconomy, error) {
	params := map[string]interface{}{
		"userID": userID,
	}

	sql := fmt.Sprintf(`
		SELECT * 
		FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" 
		WHERE user_id = :userID`, userEconomySpace())

	var res []*userEconomy
	if err := s.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to get %q record for userID:%v", userEconomySpace(), userID)
	}

	if len(res) == 0 {
		return nil, errors.Wrapf(storage.ErrNotFound,
			"unable to find %q record for userID:%v", userEconomySpace(), userID)
	}

	return res[0], nil
}

func (s *usersSource) updateUserEconomy(ue *userEconomy) error {
	nowT := uint64(time.Now().UTC().UnixNano())
	params := map[string]interface{}{
		"userId":            ue.UserID,
		"profilePictureUrl": ue.ProfilePictureURL,
		"updatedAt":         nowT,
	}

	sql := fmt.Sprintf("UPDATE %[1]v SET"+
		"profile_picture_url = :profilePictureUrl, "+
		"updated_at = :updatedAt "+
		"WHERE user_id = :userId", userEconomySpace())

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, params)),
		"failed to update user economy record for user.ID:%v", ue.UserID)
}

func (s *usersSource) createUserEconomy(u *user) error {
	nowT := uint64(time.Now().UTC().UnixNano())
	params := map[string]interface{}{
		"userId":              u.ID,
		"profilePictureUrl":   u.ProfilePictureURL,
		"balance":             0.0,
		"stakingPercentage":   0.0,
		"hashCode":            u.HashCode,
		"lastMiningStartedAt": 0,
		"stakingYears":        0,
		"createdAt":           nowT,
		"updatedAt":           nowT,
		"balanceUpdatedAt":    0,
	}

	sql := fmt.Sprintf("INSERT INTO %[1]v "+
		"(USER_ID, PROFILE_PICTURE_URL, BALANCE, STAKING_PERCENTAGE, HASH_CODE, "+
		"LAST_MINING_STARTED_AT, STAKING_YEARS, CREATED_AT, UPDATED_AT, BALANCE_UPDATED_AT) "+
		"VALUES "+
		"(:userId, :profilePictureUrl, :balance, :stakingPercentage, :hashCode, "+
		":lastMiningStartedAt, :stakingYears, :createdAt, :updatedAt, :balanceUpdatedAt)", userEconomySpace())

	return errors.Wrapf(storage.CheckSQLDMLErr(s.db.PrepareExecute(sql, params)), "failed to insert user economy record for user.ID:%v", u.ID)
}

func userEconomySpace() string {
	return "USER_ECONOMY"
}

func totalUsersSpace() string {
	return "TOTAL_USERS"
}
