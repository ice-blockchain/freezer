// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func (r *repository) GetCoinDistributionsForReview(ctx context.Context, arg *GetCoinDistributionsForReviewArg) (updCursor uint64, distributions []*PendingReview, err error) { //nolint:lll // .
	conditions, whereArgs := arg.where()
	sql := fmt.Sprintf(`SELECT * 
						FROM coin_distributions_pending_review 
						WHERE internal_id > $1 
						  AND %[1]v
						ORDER BY %[2]v 
						LIMIT $2`, strings.Join(append(conditions, "1=1"), " AND "), strings.Join(append(arg.orderBy(), "internal_id asc"), ", "))
	result, err := storage.Select[coinDistribution](ctx, r.db, sql, append([]any{arg.Cursor, arg.Limit}, whereArgs...)...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "failed to select coin_distributions_pending_review for %#v", arg)
	}
	updCursor = 0
	if uint64(len(result)) == arg.Limit {
		updCursor = result[len(result)-1].InternalID
	}
	distributions = make([]*PendingReview, len(result)) //nolint:makezero // .
	for i, d := range result {
		d.PendingReview.Ice = float64(d.PendingReview.IceInternal) / 100
		distributions[i] = d.PendingReview
	}

	return
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
	if a.ReferredByUsernameKeyword != "" {
		a.ReferredByUsernameKeyword = strings.ReplaceAll(a.ReferredByUsernameKeyword, "!", "!!")
		a.ReferredByUsernameKeyword = strings.ReplaceAll(a.ReferredByUsernameKeyword, "%", "!%")
		a.ReferredByUsernameKeyword = strings.ReplaceAll(a.ReferredByUsernameKeyword, "_", "!_")
		a.ReferredByUsernameKeyword = strings.ReplaceAll(a.ReferredByUsernameKeyword, "[", "![")
		a.ReferredByUsernameKeyword = a.ReferredByUsernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("referred_by_username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(a.ReferredByUsernameKeyword))
		i++
	}
	if a.UsernameKeyword != "" {
		a.UsernameKeyword = strings.ReplaceAll(a.UsernameKeyword, "!", "!!")
		a.UsernameKeyword = strings.ReplaceAll(a.UsernameKeyword, "%", "!%")
		a.UsernameKeyword = strings.ReplaceAll(a.UsernameKeyword, "_", "!_")
		a.UsernameKeyword = strings.ReplaceAll(a.UsernameKeyword, "[", "![")
		a.UsernameKeyword = a.UsernameKeyword + "%"
		conditions = append(conditions, fmt.Sprintf("username LIKE $%v ESCAPE '!'", i))
		args = append(args, strings.ToLower(a.UsernameKeyword))
	}

	return conditions, args
}
