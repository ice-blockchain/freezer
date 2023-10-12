// SPDX-License-Identifier: ice License 1.0

package main

import (
	"context"
	_ "embed"
	"fmt"
	"sort"
	"strings"
	"sync"
	stdlibtime "time"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/freezer/tokenomics"
	"github.com/ice-blockchain/santa/badges"
	levelsandroles "github.com/ice-blockchain/santa/levels-and-roles"
	"github.com/ice-blockchain/santa/tasks"
	appCfg "github.com/ice-blockchain/wintr/config"
	storagePG "github.com/ice-blockchain/wintr/connectors/storage/v2"
	storage "github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

const (
	applicationYamlUsersKey      = "users"
	applicationYamlKeyMiner      = "miner"
	applicationYamlKeyTokenomics = "tokenomics"
	applicationYamlKeySanta      = "santa"

	requestDeadline = 30 * stdlibtime.Second
)

var (
	//go:embed DDLSanta.sql
	ddlSanta string

	//go:embed DDLEskimo.sql
	ddlUsers string

	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfgTokenomics configTokenomics

	//nolint:gochecknoglobals // Singleton & global config mounted only during bootstrap.
	cfgSanta configSanta
)

type (
	RecalculatedAchievementsAtField struct {
		RecalculatedAchievementsAt *time.Time `redis:"recalculated_achievements_at,omitempty"`
	}

	configTokenomics struct {
		tokenomics.Config `mapstructure:",squash"` //nolint:tagliatelle // Nope.
		Workers           int64                    `yaml:"workers"`
		BatchSize         int64                    `yaml:"batchSize"`
		Development       bool                     `yaml:"development"`
	}
	configSanta struct {
		Milestones                               map[badges.Type]badges.AchievingRange `yaml:"milestones"`
		RequiredInvitedFriendsToBecomeAmbassador uint64                                `yaml:"requiredInvitedFriendsToBecomeAmbassador"`
		RequiredFriendsInvited                   uint64                                `yaml:"requiredFriendsInvited"`
	}
	updater struct {
		db       storage.DB
		dbSanta  *storagePG.DB
		dbEskimo *storagePG.DB
		wg       *sync.WaitGroup
	}
	updatedUser struct {
		Badges         *badgesUser
		Tasks          *taskUser
		LevelsAndRoles *levelsAndRolesUser
		FriendsInvited *friendsInvitedUser
		UserID         string
		model.DeserializedUsersKey
		ActualBalance        int64
		ActualFriendsInvited uint64
		CompletedTasks       bool
	}
	user struct {
		model.DeserializedUsersKey
		model.UserIDField
		model.BalanceTotalStandardField
		model.BalanceTotalPreStakingField
		RecalculatedAchievementsAtField
	}
	recalculatedUserUpdated struct {
		model.DeserializedRecalculatedUsersKey
		RecalculatedAchievementsAtField
	}
	eskimoUser struct {
		UserID         string `json:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2" db:"user_id"`
		FriendsInvited uint64 `json:"friendsInvited" example:"22" db:"friends_invited"`
	}
	badgesUser struct {
		AchievedBadges  *users.Enum[badges.Type] `json:"achievedBadges,omitempty" example:"c1,l1,l2,c2"`
		UserID          string                   `json:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2" db:"user_id"`
		FriendsInvited  int64                    `json:"friendsInvited,omitempty" example:"3"`
		Balance         int64                    `json:"balance,omitempty" example:"1232323232"`
		CompletedLevels int64                    `json:"completedLevels,omitempty" example:"3"`
	}
	taskUser struct {
		CompletedTasks       *users.Enum[tasks.Type] `json:"completedTasks,omitempty" example:"claim_username,start_mining"`
		PseudoCompletedTasks *users.Enum[tasks.Type] `json:"pseudoCompletedTasks,omitempty" example:"claim_username,start_mining"`
		UserID               string                  `json:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2" db:"user_id"`
		FriendsInvited       uint64                  `json:"friendsInvited,omitempty" example:"3"`
	}
	levelsAndRolesUser struct {
		EnabledRoles   *users.Enum[levelsandroles.RoleType]
		UserID         string `json:"userId" example:"did:ethr:0x4B73C58370AEfcEf86A6021afCDe5673511376B2" db:"user_id"`
		FriendsInvited uint64 `json:"friendsInvited,omitempty" example:"3"`
		CompletedTasks uint64 `json:"completedTasks,omitempty" example:"3"`
	}
	friendsInvitedUser struct {
		UserID         string `json:"userId" example:"edfd8c02-75e0-4687-9ac2-1ce4723865c4"`
		FriendsInvited uint64 `json:"friendsInvited" example:"5" db:"invited_count"`
	}
)

func main() {
	appCfg.MustLoadFromKey(applicationYamlKeyTokenomics, &cfgTokenomics.Config)
	appCfg.MustLoadFromKey(applicationYamlKeyMiner, &cfgTokenomics)
	appCfg.MustLoadFromKey(applicationYamlKeySanta, &cfgSanta)

	dbEskimo := storagePG.MustConnect(context.Background(), ddlUsers, applicationYamlUsersKey)
	dbSanta := storagePG.MustConnect(context.Background(), ddlSanta, applicationYamlKeySanta)

	if err := dbEskimo.Ping(context.Background()); err != nil {
		log.Panic("can't ping users db", err)
	}
	if err := dbSanta.Ping(context.Background()); err != nil {
		log.Panic("can't ping santa db", err)
	}
	upd := &updater{
		db:       storage.MustConnect(context.Background(), applicationYamlKeyTokenomics, int(cfgTokenomics.Workers)),
		dbSanta:  dbSanta,
		dbEskimo: dbEskimo,
		wg:       new(sync.WaitGroup),
	}
	defer upd.db.Close()
	defer upd.dbEskimo.Close()
	defer upd.dbSanta.Close()
	if resp := upd.db.Ping(context.Background()); resp.Err() != nil || resp.Val() != "PONG" {
		if resp.Err() == nil {
			resp.SetErr(errors.Errorf("response `%v` is not `PONG`", resp.Val()))
		}

		log.Panic(errors.Wrap(resp.Err(), "failed to ping DB"))
	}
	if !upd.db.IsRW(context.Background()) {
		log.Panic("db is not writeable")
	}
	upd.wg.Add(int(cfgTokenomics.Workers))
	for workerNumber := int64(0); workerNumber < cfgTokenomics.Workers; workerNumber++ {
		go func(wn int64) {
			defer upd.wg.Done()
			upd.update(context.Background(), wn)
		}(workerNumber)
	}
	upd.wg.Wait()
}

func (u *updater) update(ctx context.Context, workerNumber int64) {
	var (
		batchNumber              int64
		workers                  = cfgTokenomics.Workers
		batchSize                = cfgTokenomics.BatchSize
		userKeys                 = make([]string, 0, batchSize)
		userResults              = make([]*user, 0, batchSize)
		userRecalculatedResults  = make([]*recalculatedUserUpdated, 0, batchSize)
		recalculatedUsersUpdated = make([]*recalculatedUserUpdated, 0, batchSize)
		updatedUsers             = make(map[string]*updatedUser, batchSize)
		errs                     = make([]error, 0, batchSize)
	)
	resetVars := func() {
		userKeys = userKeys[:0]
		userResults = userResults[:0]
		userRecalculatedResults = userRecalculatedResults[:0]
		recalculatedUsersUpdated = recalculatedUsersUpdated[:0]
		errs = errs[:0]
	}

	for ctx.Err() == nil {
		/******************************************************************************************************************************************************
			1. Fetching a new batch of users from redis.
		******************************************************************************************************************************************************/
		if len(userKeys) == 0 {
			for ix := batchNumber * batchSize; ix < (batchNumber+1)*batchSize; ix++ {
				userKeys = append(userKeys, model.SerializedUsersKey((workers*ix)+workerNumber))
			}
		}
		reqCtx, reqCancel := context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[user](reqCtx, u.db, userKeys, &userResults); err != nil {
			log.Panic(errors.Wrapf(err, "[miner] failed to get users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		if len(userResults) == 0 {
			reqCancel()
			log.Info("updating finished for worker:", workerNumber)

			break
		}
		userKeys = userKeys[:0]
		if len(userKeys) == 0 {
			for ix := batchNumber * batchSize; ix < (batchNumber+1)*batchSize; ix++ {
				userKeys = append(userKeys, model.SerializedRecalculatedUsersKey((workers*ix)+workerNumber))
			}
		}
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if err := storage.Bind[recalculatedUserUpdated](reqCtx, u.db, userKeys, &userRecalculatedResults); err != nil {
			log.Panic(errors.Wrapf(err, "[miner] failed to get users for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
			reqCancel()

			continue
		}
		reqCancel()

		/******************************************************************************************************************************************************
			2. Fetching friends invited count, completed tasks and achieved badges for them.
		******************************************************************************************************************************************************/
		var userIDs []string
	outer:
		for _, usr := range userResults {
			if usr.UserID == "" {
				continue
			}
			for _, recalculatedUser := range userRecalculatedResults {
				if recalculatedUser.ID == usr.ID && !recalculatedUser.RecalculatedAchievementsAt.IsNil() {
					continue outer
				}
			}
			updatedUsers[usr.UserID] = &updatedUser{
				DeserializedUsersKey: usr.DeserializedUsersKey,
				UserID:               usr.UserID,
				ActualBalance:        int64(usr.BalanceTotalStandard + usr.BalanceTotalPreStaking),
			}
			userIDs = append(userIDs, usr.UserID)
		}
		if updatedUsers == nil || len(userIDs) == 0 {
			log.Debug("no user ids, worker:", workerNumber)

			return
		}
		u.fetchActualFriendsInvited(ctx, userIDs, updatedUsers)
		u.fetchTasks(ctx, userIDs, updatedUsers)
		u.fetchAchievedBadges(ctx, userIDs, updatedUsers)
		u.fetchLevelsAndRoles(ctx, userIDs, updatedUsers)
		u.fetchFriendsInvited(ctx, userIDs, updatedUsers)

		/******************************************************************************************************************************************************
			3. Updating santa.
		******************************************************************************************************************************************************/
		for _, usr := range updatedUsers {
			if err := u.updateBadgesAndStatistics(ctx, usr); err != nil {
				log.Panic("can't update badges and badges statistics, userID:", usr.UserID)
			}
			if err := u.updateLevelsAndRoles(ctx, usr); err != nil {
				log.Panic("can't update levels and roles, userID:", usr.UserID)
			}
			if err := u.updateTasks(ctx, usr); err != nil {
				log.Panic("can't update tasks, userID:", usr.UserID)
			}
			if err := u.updateFriendsInvited(ctx, usr); err != nil {
				log.Panic("can't update friends invited, userID:", usr.UserID)
			}
			recalculatedUsersUpdated = append(recalculatedUsersUpdated, &recalculatedUserUpdated{
				DeserializedRecalculatedUsersKey: model.DeserializedRecalculatedUsersKey{ID: usr.ID},
				RecalculatedAchievementsAtField:  RecalculatedAchievementsAtField{RecalculatedAchievementsAt: time.Now()},
			})
		}

		/******************************************************************************************************************************************************
			4. Persisting recalculated status.
		******************************************************************************************************************************************************/
		pipeliner := u.db.Pipeline()
		reqCtx, reqCancel = context.WithTimeout(context.Background(), requestDeadline)
		if responses, err := pipeliner.Pipelined(reqCtx, func(pipeliner redis.Pipeliner) error {
			for _, value := range recalculatedUsersUpdated {
				if err := pipeliner.HSet(reqCtx, value.Key(), storage.SerializeValue(value)...).Err(); err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			log.Panic(errors.Wrapf(err, "[updater] [1]failed to persist update process for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
		} else {
			for _, response := range responses {
				if err = response.Err(); err != nil {
					errs = append(errs, errors.Wrapf(err, "failed to `%v`", response.FullName()))
				}
			}
			if err = multierror.Append(nil, errs...).ErrorOrNil(); err != nil {
				log.Error(errors.Wrapf(err, "[updater] [2]failed to persist update progress for batchNumber:%v,workerNumber:%v", batchNumber, workerNumber))
				reqCancel()
				resetVars()

				continue
			}
		}

		reqCancel()
		batchNumber++
		resetVars()
	}
}

func (u *updater) fetchActualFriendsInvited(ctx context.Context, userIDs []string, updatedUsers map[string]*updatedUser) {
	sql := `SELECT 
				u.id 		 AS user_id,
				COUNT(t1.id) AS friends_invited
			FROM users u
			LEFT JOIN users t1
				ON t1.referred_by = u.id
				   AND u.username != u.id
			WHERE u.id = ANY($1)
				  AND u.referred_by != u.id
				  AND u.username != u.id
			GROUP BY u.id`
	result, err := storagePG.Select[eskimoUser](ctx, u.dbEskimo, sql, userIDs)
	if err != nil {
		log.Panic("error on trying to get friends invited count", userIDs, err)
	}
	if len(result) == 0 {
		log.Debug("no results for: ", userIDs, err)

		return
	}
	for _, r := range result {
		updatedUsers[r.UserID].ActualFriendsInvited = r.FriendsInvited
	}
}

func (u *updater) fetchTasks(ctx context.Context, userIDs []string, updatedUsers map[string]*updatedUser) {
	sql := `SELECT 
				user_id,
				friends_invited,
				completed_tasks,
				pseudo_completed_tasks
			FROM task_progress
				WHERE user_id = ANY($1)`

	result, err := storagePG.Select[taskUser](ctx, u.dbSanta, sql, userIDs)
	if err != nil {
		log.Panic("error on trying to get tasks ", userIDs, err)
	}
	if len(result) == 0 {
		log.Debug("no results for: ", userIDs, err)

		return
	}
	for _, r := range result {
		updatedUsers[r.UserID].Tasks = &taskUser{
			UserID:         r.UserID,
			FriendsInvited: r.FriendsInvited,
		}
		if r.CompletedTasks != nil {
			for _, task := range *r.CompletedTasks {
				if task == tasks.InviteFriendsType {
					updatedUsers[r.UserID].CompletedTasks = true
					break
				}
			}
		}
		if r.PseudoCompletedTasks != nil && !updatedUsers[r.UserID].CompletedTasks {
			for _, task := range *r.PseudoCompletedTasks {
				if task == tasks.InviteFriendsType {
					updatedUsers[r.UserID].CompletedTasks = true
					break
				}
			}
		}
	}
}

func (u *updater) fetchFriendsInvited(ctx context.Context, userIDs []string, updatedUsers map[string]*updatedUser) {
	sql := `SELECT 
				user_id,
				invited_count
			FROM friends_invited
			WHERE user_id = ANY($1)`

	result, err := storagePG.Select[friendsInvitedUser](ctx, u.dbSanta, sql, userIDs)
	if err != nil {
		log.Panic("error on trying to get tasks ", userIDs, err)
	}
	if len(result) == 0 {
		log.Debug("no results for: ", userIDs, err)

		return
	}
	for _, r := range result {
		updatedUsers[r.UserID].FriendsInvited = &friendsInvitedUser{
			UserID:         r.UserID,
			FriendsInvited: r.FriendsInvited,
		}
	}
}

func (u *updater) fetchAchievedBadges(ctx context.Context, userIDs []string, updatedUsers map[string]*updatedUser) {
	sql := `SELECT 
				user_id,
				achieved_badges,
				completed_levels,
				friends_invited,
				balance
			FROM badge_progress
				WHERE user_id = ANY($1)`

	result, err := storagePG.Select[badgesUser](ctx, u.dbSanta, sql, userIDs)
	if err != nil {
		log.Panic("error on trying to get friends invited counts ", userIDs, err)
	}
	if len(result) == 0 {
		log.Debug("no results for: ", userIDs, err)

		return
	}
	for _, r := range result {
		updatedUsers[r.UserID].Badges = &badgesUser{
			AchievedBadges:  r.AchievedBadges,
			UserID:          r.UserID,
			FriendsInvited:  r.FriendsInvited,
			Balance:         r.Balance,
			CompletedLevels: r.CompletedLevels,
		}
	}
}

func (u *updater) fetchLevelsAndRoles(ctx context.Context, userIDs []string, updatedUsers map[string]*updatedUser) {
	sql := `SELECT 
				user_id,
				friends_invited,
				enabled_roles,
				completed_tasks
			FROM levels_and_roles_progress
				WHERE user_id = ANY($1)`

	result, err := storagePG.Select[levelsAndRolesUser](ctx, u.dbSanta, sql, userIDs)
	if err != nil {
		log.Panic("error on trying to get friends invited counts ", userIDs, err)
	}
	if len(result) == 0 {
		log.Debug("no results for: ", userIDs, err)

		return
	}
	for _, r := range result {
		updatedUsers[r.UserID].LevelsAndRoles = &levelsAndRolesUser{
			UserID:         r.UserID,
			FriendsInvited: r.FriendsInvited,
			EnabledRoles:   r.EnabledRoles,
			CompletedTasks: r.CompletedTasks,
		}
	}
}

func (u *updater) updateBadgesAndStatistics(ctx context.Context, usr *updatedUser) error {
	achievedBadges, newBadgesTypeCount := u.reEvaluateEnabledBadges(usr.Badges.AchievedBadges, usr.ActualFriendsInvited, usr.ActualBalance)
	var completedLevelsSQL string
	if usr.CompletedTasks && usr.ActualFriendsInvited < cfgSanta.RequiredFriendsInvited {
		completedLevelsSQL = ",completed_levels = GREATEST(completed_levels - 1, 0)"
	}
	newBadgesTypeCount = u.diffBadgeStatistics(usr, newBadgesTypeCount)
	var achievedBadgesChanged bool
	for _, count := range newBadgesTypeCount {
		if count != 0 {
			achievedBadgesChanged = true

			break
		}
	}
	if achievedBadgesChanged || completedLevelsSQL != "" || usr.ActualFriendsInvited != uint64(usr.Badges.FriendsInvited) || usr.Badges.Balance != usr.ActualBalance {
		sql := fmt.Sprintf(`UPDATE badge_progress
								SET friends_invited = $2,
									balance = $3,
									achieved_badges = $4
									%v
							WHERE user_id = $1`, completedLevelsSQL)
		if _, err := storagePG.Exec(ctx, u.dbSanta, sql, usr.UserID, usr.ActualFriendsInvited, usr.ActualBalance, achievedBadges); err != nil {
			return errors.Wrapf(err, "failed to set badge_progress.friends_invited, userID:%v, friendsInvited:%v, balance:%v", usr.UserID, usr.ActualFriendsInvited, usr.ActualBalance)
		}
		var mErr *multierror.Error
		for badgeType, val := range newBadgesTypeCount {
			if val != 0 {
				sign := "+"
				if val < 0 {
					sign = "-"
					val *= -1
				}
				sql := fmt.Sprintf(`UPDATE badge_statistics
										SET achieved_by = GREATEST(achieved_by %v $1, 0)
									WHERE badge_type = $2`, sign)
				_, err := storagePG.Exec(ctx, u.dbSanta, sql, val, badgeType)
				mErr = multierror.Append(errors.Wrapf(err, "failed to update badge_statistics, userID:%v, badgeType:%v, val:%v", usr.UserID, badgeType, val))
			}
		}

		return multierror.Append(mErr, nil).ErrorOrNil()
	}

	return nil
}

func (u *updater) diffBadgeStatistics(usr *updatedUser, newBadgesTypeCount map[badges.Type]int64) map[badges.Type]int64 {
	oldBadgesTypeCounts := make(map[badges.Type]int64, len(badges.AllTypes))
	oldGroupCounts := make(map[badges.GroupType]int64, len(badges.AllGroups))
	for _, badge := range *usr.Badges.AchievedBadges {
		switch badges.GroupTypeForEachType[badge] {
		case badges.CoinGroupType:
			oldBadgesTypeCounts[badge]++
			oldGroupCounts[badges.CoinGroupType]++
		case badges.SocialGroupType:
			oldBadgesTypeCounts[badge]++
			oldGroupCounts[badges.SocialGroupType]++
		}
	}
	if newBadgesTypeCount != nil {
		for _, key := range badges.AllTypes {
			if _, ok1 := oldBadgesTypeCounts[key]; ok1 {
				if _, ok2 := newBadgesTypeCount[key]; ok2 {
					newBadgesTypeCount[key] = newBadgesTypeCount[key] - oldBadgesTypeCounts[key]
				} else {
					newBadgesTypeCount[key] -= oldBadgesTypeCounts[key]
				}
			}
		}
	}

	return newBadgesTypeCount
}

func (u *updater) updateLevelsAndRoles(ctx context.Context, usr *updatedUser) error {
	enabledRoles := u.reEvaluateEnabledRole(usr.ActualFriendsInvited)
	var roleChanged bool
	if (usr.LevelsAndRoles.EnabledRoles != nil && enabledRoles == nil) ||
		(usr.LevelsAndRoles.EnabledRoles == nil && enabledRoles != nil) {
		roleChanged = true
	}
	var completedTasksSQL, completedLevelsSQL string
	if usr.CompletedTasks && usr.ActualFriendsInvited < cfgSanta.RequiredFriendsInvited {
		completedTasksSQL = ", completed_tasks = GREATEST(completed_tasks - 1, 0)"
		completedLevelsSQL = ",completed_levels = array_remove(completed_levels, '11')" // We know for sure from config file this level id that need to be removed.
	}
	if completedTasksSQL != "" || completedLevelsSQL != "" || usr.ActualFriendsInvited != usr.LevelsAndRoles.FriendsInvited || roleChanged {
		sql := fmt.Sprintf(`UPDATE levels_and_roles_progress
								SET friends_invited = $2,
									enabled_roles = $3
									%v
									%v
							WHERE user_id = $1`, completedTasksSQL, completedLevelsSQL)
		_, err := storagePG.Exec(ctx, u.dbSanta, sql, usr.UserID, usr.ActualFriendsInvited, enabledRoles)

		return errors.Wrapf(err, "failed to set levels_and_roles_progress.friends_invited, userID:%v, friendsInvited:%v", usr.UserID, usr.ActualFriendsInvited)
	}

	return nil
}

func (u *updater) updateTasks(ctx context.Context, usr *updatedUser) error {
	var completedTasksSQL string
	if usr.CompletedTasks && usr.ActualFriendsInvited < cfgSanta.RequiredFriendsInvited {
		completedTasksSQL = `, completed_tasks = array_remove(completed_tasks, 'invite_friends')
							 , pseudo_completed_tasks = array_remove(pseudo_completed_tasks, 'invite_friends')`
	}
	if completedTasksSQL != "" || (usr.Tasks != nil && usr.ActualFriendsInvited != usr.Tasks.FriendsInvited) {
		sql := fmt.Sprintf(`UPDATE task_progress
								SET friends_invited = $2
								%v
							WHERE user_id = $1`, completedTasksSQL)
		_, err := storagePG.Exec(ctx, u.dbSanta, sql, usr.UserID, usr.ActualFriendsInvited)

		return errors.Wrapf(err, "failed to set task_progress.friends_invited, userID:%v, friendsInvited:%v", usr.UserID, usr.ActualFriendsInvited)
	}

	return nil
}

func (u *updater) updateFriendsInvited(ctx context.Context, usr *updatedUser) error {
	if usr.ActualFriendsInvited != usr.FriendsInvited.FriendsInvited {
		sql := `INSERT INTO friends_invited(user_id, invited_count) VALUES ($1, $2)
				ON CONFLICT(user_id) DO UPDATE SET
					invited_count = EXCLUDED.invited_count
				WHERE friends_invited.invited_count != EXCLUDED.invited_count`
		_, err := storagePG.Exec(ctx, u.dbSanta, sql, usr.UserID, usr.ActualFriendsInvited)

		return errors.Wrapf(err, "failed to set task_progress.friends_invited, userID:%v, friendsInvited:%v", usr.UserID, usr.ActualFriendsInvited)
	}

	return nil
}

func (u *updater) reEvaluateEnabledRole(friendsInvited uint64) *users.Enum[levelsandroles.RoleType] {
	if friendsInvited >= cfgSanta.RequiredInvitedFriendsToBecomeAmbassador {
		completedLevels := append(make(users.Enum[levelsandroles.RoleType], 0, len(&levelsandroles.AllRoleTypesThatCanBeEnabled)), levelsandroles.AmbassadorRoleType)

		return &completedLevels
	}

	return nil
}

func (u *updater) reEvaluateEnabledBadges(alreadyAchievedBadges *users.Enum[badges.Type], friendsInvited uint64, balance int64) (achievedBadges users.Enum[badges.Type], badgesTypeCounts map[badges.Type]int64) {
	badgesTypeCounts = make(map[badges.Type]int64)
	achievedBadges = make(users.Enum[badges.Type], 0, len(&badges.AllTypes))
	if alreadyAchievedBadges != nil {
		for _, badge := range *alreadyAchievedBadges {
			if strings.HasPrefix(string(badge), "l") {
				achievedBadges = append(achievedBadges, badge)
			}
		}
	}
	for _, badgeType := range &badges.AllTypes {
		var achieved bool
		switch badges.GroupTypeForEachType[badgeType] {
		case badges.CoinGroupType:
			if balance > 0 {
				achieved = uint64(balance) >= cfgSanta.Milestones[badgeType].FromInclusive
			}
		case badges.SocialGroupType:
			achieved = uint64(friendsInvited) >= cfgSanta.Milestones[badgeType].FromInclusive
		}
		if achieved {
			achievedBadges = append(achievedBadges, badgeType)
			badgesTypeCounts[badgeType]++
		}
	}
	if len(achievedBadges) == 0 {
		return nil, nil
	}
	sort.SliceStable(achievedBadges, func(i, j int) bool {
		return badges.AllTypeOrder[achievedBadges[i]] < badges.AllTypeOrder[achievedBadges[j]]
	})

	return achievedBadges, badgesTypeCounts
}
