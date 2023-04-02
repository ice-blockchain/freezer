// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/coin"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func (s *usersTableSource) Process(ctx context.Context, msg *messagebroker.Message) error { //nolint:gocognit // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(msg.Value) == 0 {
		return nil
	}
	var usr users.UserSnapshot
	if err := json.UnmarshalContext(ctx, msg.Value, &usr); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &usr)
	}
	if (usr.User == nil || usr.User.ID == "") && (usr.Before == nil || usr.Before.ID == "") {
		return nil
	}

	if usr.User == nil || usr.User.ID == "" {
		return errors.Wrapf(s.deleteUser(ctx, usr.Before), "failed to delete user:%#v", usr.Before)
	}

	if err := s.replaceUser(ctx, usr.User); err != nil {
		return errors.Wrapf(err, "failed to replace user:%#v", usr.User)
	}

	return nil
}

func (s *usersTableSource) deleteUser(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	if err := s.removeBalanceFromT0AndTMinus1(ctx, usr); err != nil {
		return errors.Wrapf(err, "failed to removeBalanceFromT0AndTMinus1 for user:%#v", usr)
	}
	_, err := storage.Exec(ctx, s.db, `DELETE FROM users WHERE user_id = $1`, usr.ID)

	return errors.Wrapf(err, "failed to delete userID:%v", usr.ID)
}

func (s *usersTableSource) removeBalanceFromT0AndTMinus1(ctx context.Context, usr *users.User) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`SELECT reverse_t0_balance.amount AS total_reverse_t0_amount,
								reverse_tminus1_balance.amount AS total_reverse_t_minus1_amount,
								negative_t0_balance.amount AS negative_reverse_t0_amount,
								negative_tminus1_balance.amount AS negative_reverse_t_minus1_amount,
				   			    t0.user_id AS t0_user_id,
								coalesce(tminus1.user_id,'') AS t_minus1_user_id
						FROM users u
							 JOIN users t0
							   ON t0.user_id = u.referred_by
							  AND t0.user_id != u.user_id
						LEFT JOIN users tminus1
							   ON tminus1.user_id = t0.referred_by
							  AND tminus1.user_id != t0.user_id
						LEFT JOIN balances_worker reverse_t0_balance
							   ON reverse_t0_balance.worker_index = $1
							  AND reverse_t0_balance.user_id = u.user_id
							  AND reverse_t0_balance.negative = FALSE
							  AND reverse_t0_balance.type = %[1]v
							  AND reverse_t0_balance.type_detail =  '%[2]v_' || t0.user_id
						LEFT JOIN balances_worker reverse_tminus1_balance
							   ON reverse_tminus1_balance.worker_index = $1
							  AND reverse_tminus1_balance.user_id = u.user_id
							  AND reverse_tminus1_balance.negative = FALSE
							  AND reverse_tminus1_balance.type = %[1]v
							  AND tminus1.user_id IS NOT NULL
							  AND reverse_tminus1_balance.type_detail =  '%[3]v_' || tminus1.user_id
						LEFT JOIN balances_worker negative_t0_balance
							   ON negative_t0_balance.worker_index = $1
							  AND negative_t0_balance.user_id = u.user_id
							  AND negative_t0_balance.negative = TRUE
							  AND negative_t0_balance.type = %[1]v
							  AND negative_t0_balance.type_detail =  '%[2]v_' || t0.user_id
						LEFT JOIN balances_worker negative_tminus1_balance
							   ON negative_tminus1_balance.worker_index = $1
							  AND negative_tminus1_balance.user_id = u.user_id
							  AND negative_tminus1_balance.negative = TRUE
							  AND negative_tminus1_balance.type = %[1]v
							  AND tminus1.user_id IS NOT NULL
							  AND negative_tminus1_balance.type_detail =  '%[3]v_' || tminus1.user_id
					    WHERE u.user_id = $2`, totalNoPreStakingBonusBalanceType, reverseT0BalanceTypeDetail, reverseTMinus1BalanceTypeDetail)
	res, err := storage.Get[struct {
		TotalReverseT0Amount, TotalReverseTMinus1Amount,
		NegativeReverseT0Amount, NegativeReverseTMinus1Amount *coin.ICEFlake
		T0UserID, TMinus1UserID string
	}](ctx, s.db, sql, int16(usr.HashCode%uint64(s.cfg.WorkerCount)), usr.ID)
	if err != nil {
		if storage.IsErr(err, storage.ErrNotFound) {
			return nil
		}

		return errors.Wrapf(err, "failed to get reverse t0 and t-1 balance information for userID:%v", usr.ID)
	}
	cmds := make([]*AddBalanceCommand, 0, 1+1+1+1)
	if !res.TotalReverseT0Amount.IsZero() {
		cmds = append(cmds, &AddBalanceCommand{
			Balances: &Balances[coin.ICEFlake]{
				T1:     res.TotalReverseT0Amount,
				UserID: res.T0UserID,
			},
			EventID: fmt.Sprintf("t1_referral_account_deletion_positive_balance_%v", usr.ID),
		})
	}
	if !res.NegativeReverseT0Amount.IsZero() {
		negative := true
		cmds = append(cmds, &AddBalanceCommand{
			Balances: &Balances[coin.ICEFlake]{
				T1:     res.NegativeReverseT0Amount,
				UserID: res.T0UserID,
			},
			EventID:  fmt.Sprintf("t1_referral_account_deletion_negative_balance_%v", usr.ID),
			Negative: &negative,
		})
	}
	if !res.TotalReverseTMinus1Amount.IsZero() {
		cmds = append(cmds, &AddBalanceCommand{
			Balances: &Balances[coin.ICEFlake]{
				T2:     res.TotalReverseTMinus1Amount,
				UserID: res.TMinus1UserID,
			},
			EventID: fmt.Sprintf("t2_referral_account_deletion_positive_balance_%v", usr.ID),
		})
	}
	if !res.NegativeReverseTMinus1Amount.IsZero() {
		negative := true
		cmds = append(cmds, &AddBalanceCommand{
			Balances: &Balances[coin.ICEFlake]{
				T2:     res.NegativeReverseTMinus1Amount,
				UserID: res.TMinus1UserID,
			},
			EventID:  fmt.Sprintf("t2_referral_account_deletion_negative_balance_%v", usr.ID),
			Negative: &negative,
		})
	}

	return errors.Wrapf(executeBatchConcurrently(ctx, s.sendAddBalanceCommandMessage, cmds), "failed to sendAddBalanceCommandMessages for %#v", cmds)
}

func (s *usersTableSource) replaceUser(ctx context.Context, usr *users.User) error { //nolint:funlen // .
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := `INSERT INTO users (created_at, updated_at, user_id, referred_by, username, first_name, last_name, profile_picture_name, mining_blockchain_account_address, blockchain_account_address, hash_code  , hide_ranking, verified)
					   VALUES ($1		 , $2		 , $3	  , $4		   , $5	     , $6		 , $7	    , $8			      , $9								 , $10					     , $11::bigint, $12			, $13)
		    ON CONFLICT (user_id)
		    	DO UPDATE
		    		  SET updated_at 						= EXCLUDED.updated_at,
		    		  	  referred_by 						= EXCLUDED.referred_by,
		    		  	  username 							= EXCLUDED.username,
		    		  	  first_name 						= EXCLUDED.first_name,
		    		  	  last_name 						= EXCLUDED.last_name,
		    		  	  profile_picture_name 				= EXCLUDED.profile_picture_name,
		    		  	  mining_blockchain_account_address = EXCLUDED.mining_blockchain_account_address,
		    		  	  blockchain_account_address 		= EXCLUDED.blockchain_account_address,
		    		  	  hide_ranking 						= EXCLUDED.hide_ranking
			  		WHERE coalesce(users.referred_by,'') 						!= coalesce(EXCLUDED.referred_by,'')
			  		   OR coalesce(users.username,'') 							!= coalesce(EXCLUDED.username,'')
			  		   OR coalesce(users.first_name,'') 						!= coalesce(EXCLUDED.first_name,'')
			  		   OR coalesce(users.last_name,'') 							!= coalesce(EXCLUDED.last_name,'')
			  		   OR coalesce(users.profile_picture_name,'') 				!= coalesce(EXCLUDED.profile_picture_name,'')
			  		   OR coalesce(users.mining_blockchain_account_address,'') 	!= coalesce(EXCLUDED.mining_blockchain_account_address,'')
			  		   OR coalesce(users.blockchain_account_address,'') 		!= coalesce(EXCLUDED.blockchain_account_address,'')
			  		   OR coalesce(users.hide_ranking,false) 					!= coalesce(EXCLUDED.hide_ranking,false)
			  		   OR coalesce(users.verified,false) 						!= coalesce(EXCLUDED.verified,false)`
	verified := false
	if usr.Verified != nil && *usr.Verified {
		verified = true
	}
	args := append(make([]any, 0, 13), //nolint:gomnd // .
		*usr.CreatedAt.Time,
		*usr.UpdatedAt.Time,
		usr.ID,
		usr.ReferredBy,
		usr.Username,
		usr.FirstName,
		usr.LastName,
		s.pictureClient.StripDownloadURL(usr.ProfilePictureURL),
		usr.MiningBlockchainAccountAddress,
		usr.BlockchainAccountAddress,
		int64(usr.HashCode),
		s.hideRanking(usr),
		verified)
	if _, err := storage.Exec(ctx, s.db, sql, args...); err != nil {
		return errors.Wrapf(err, "failed to replace user:%#v", usr)
	}

	return multierror.Append(
		errors.Wrapf(s.updateBlockchainBalanceSynchronizationWorkerBlockchainAccountAddress(ctx, usr), "failed to updateBlockchainBalanceSynchronizationWorkerBlockchainAccountAddress for usr:%#v", usr), //nolint:lll // .
		errors.Wrapf(s.initializeBalanceRecalculationWorker(ctx, usr), "failed to initializeBalanceRecalculationWorker for %#v", usr),
		errors.Wrapf(s.initializeMiningRatesRecalculationWorker(ctx, usr), "failed to initializeMiningRatesRecalculationWorker for %#v", usr),
		errors.Wrapf(s.initializeBlockchainBalanceSynchronizationWorker(ctx, usr), "failed to initializeBlockchainBalanceSynchronizationWorker for %#v", usr),
		errors.Wrapf(s.initializeExtraBonusProcessingWorker(ctx, usr), "failed to initializeExtraBonusProcessingWorker for %#v", usr),
		errors.Wrapf(s.awardRegistrationICECoinsBonus(ctx, usr), "failed to awardRegistrationICECoinsBonus for %#v", usr))
}

func (*usersTableSource) hideRanking(usr *users.User) (hideRanking bool) {
	if usr.HiddenProfileElements != nil {
		for _, element := range *usr.HiddenProfileElements {
			if users.GlobalRankHiddenProfileElement == element {
				hideRanking = true

				break
			}
		}
	}

	return hideRanking
}

func (s *usersTableSource) awardRegistrationICECoinsBonus(ctx context.Context, usr *users.User) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	cmd := &AddBalanceCommand{
		Balances: &Balances[coin.ICEFlake]{
			Total:  coin.NewAmountUint64(registrationICEFlakeBonusAmount),
			UserID: usr.ID,
		},
		EventID: "registration_ice_bonus",
	}

	return errors.Wrapf(s.sendAddBalanceCommandMessage(ctx, cmd), "failed to sendAddBalanceCommandMessage for %#v", cmd)
}
