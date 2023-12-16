// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
)

//nolint:funlen // .
func (r *repository) GetCoinDistributionsForReview(ctx context.Context, arg *GetCoinDistributionsForReviewArg) (*CoinDistributionsForReview, error) { //nolint:lll // .
	conditions, whereArgs := arg.where()
	sql := fmt.Sprintf(`SELECT * 
						FROM coin_distributions_pending_review 
						WHERE 1=1 
						  AND %[1]v
						ORDER BY %[2]v 
						LIMIT $2 OFFSET $1`, strings.Join(append(conditions, "1=1"), " AND "), strings.Join(append(arg.orderBy(), "internal_id asc"), ", "))
	result, err := storage.Select[struct {
		*PendingReview
		Day        stdlibtime.Time
		InternalID uint64
	}](ctx, r.db, sql, append([]any{arg.Cursor, arg.Limit}, whereArgs...)...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select coin_distributions_pending_review for %#v", arg)
	}
	distributions := make([]*PendingReview, len(result)) //nolint:makezero // .
	for i, d := range result {
		d.PendingReview.Ice = float64(d.PendingReview.IceInternal) / 100
		distributions[i] = d.PendingReview
	}
	conditions, whereArgs = arg.totalsWhere()
	sql = fmt.Sprintf(`SELECT count(1) AS rows,
							   coalesce(sum(ice),0) AS ice 
					   FROM coin_distributions_pending_review 
					   WHERE 1=1
						 AND %[1]v`, strings.Join(append(conditions, "1=1"), " AND "))
	total, err := storage.Get[struct {
		Rows uint64
		Ice  uint64
	}](ctx, r.db, sql, whereArgs...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to select coin_distributions_pending_review totals for %#v", arg)
	}
	nextCursor := uint64(0)
	if len(result) == int(arg.Limit) {
		nextCursor = arg.Cursor + arg.Limit
	}

	return &CoinDistributionsForReview{
		Distributions: distributions,
		Cursor:        nextCursor,
		TotalRows:     total.Rows,
		TotalIce:      float64(total.Ice) / 100,
	}, nil
}

func (a *GetCoinDistributionsForReviewArg) orderBy() []string {
	res := make([]string, 0, 4)

	if a.IceOrderBy != "" {
		res = append(res, fmt.Sprintf("ice %v", a.IceOrderBy))
	}
	if a.ReferredByUsernameOrderBy != "" {
		res = append(res, fmt.Sprintf("referred_by_username %v", a.ReferredByUsernameOrderBy))
	}
	if a.CreatedAtOrderBy != "" {
		res = append(res, fmt.Sprintf("created_at %v", a.CreatedAtOrderBy))
	}
	if a.UsernameOrderBy != "" {
		res = append(res, fmt.Sprintf("username %v", a.UsernameOrderBy))
	}

	return res
}

func (a *GetCoinDistributionsForReviewArg) where() ([]string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	i := 3
	if referredByUsernameKeyword := a.ReferredByUsernameKeyword; referredByUsernameKeyword != "" {
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "!", "!!")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "%", "!%")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "_", "!_")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "[", "![")
		referredByUsernameKeyword = referredByUsernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("referred_by_username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(referredByUsernameKeyword))
		i++
	}
	if usernameKeyword := a.UsernameKeyword; usernameKeyword != "" {
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "!", "!!")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "%", "!%")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "_", "!_")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "[", "![")
		usernameKeyword = usernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(usernameKeyword))
	}

	return conditions, args
}

func (a *GetCoinDistributionsForReviewArg) totalsWhere() ([]string, []any) {
	conditions := make([]string, 0, 2)
	args := make([]any, 0, 2)

	i := 1
	if referredByUsernameKeyword := a.ReferredByUsernameKeyword; referredByUsernameKeyword != "" {
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "!", "!!")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "%", "!%")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "_", "!_")
		referredByUsernameKeyword = strings.ReplaceAll(referredByUsernameKeyword, "[", "![")
		referredByUsernameKeyword = referredByUsernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("referred_by_username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(referredByUsernameKeyword))
		i++
	}
	if usernameKeyword := a.UsernameKeyword; usernameKeyword != "" {
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "!", "!!")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "%", "!%")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "_", "!_")
		usernameKeyword = strings.ReplaceAll(usernameKeyword, "[", "![")
		usernameKeyword = usernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(usernameKeyword))
	}

	return conditions, args
}

func (r *repository) ReviewCoinDistributions(ctx context.Context, reviewerUserID string, decision string) error {
	log.Info(fmt.Sprintf("ReviewCoinDistributions(userID:`%v`, decision:`%v`)", reviewerUserID, decision))

	return ctx.Err()
}
