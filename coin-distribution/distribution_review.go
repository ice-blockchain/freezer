// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func (r *repository) GetCoinDistributionsForReview(ctx context.Context, cursor, limit uint64) (updCursor uint64, distributions []*CoinDistibutionForReview, err error) {
	sql := `SELECT * FROM coin_distributions_pending_review WHERE internal_id > $1 ORDER BY internal_id LIMIT $2`
	args := []any{cursor, limit}
	result, err := storage.Select[coinDistribution](ctx, r.db, sql, args...)
	if err != nil {
		return 0, nil, errors.Wrapf(err, "failed to select %v records with distributions under review for > %v ", limit, cursor)
	}
	updCursor = 0
	if uint64(len(result)) == limit {
		updCursor = result[len(result)-1].InternalID
	}
	distributions = make([]*CoinDistibutionForReview, len(result)) //nolint:makezero // .
	for i, d := range result {
		distributions[i] = d.CoinDistibutionForReview
	}

	return
}
