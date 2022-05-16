// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

func (e *economy) GetTopMiners(ctx context.Context, limit, offset uint64) ([]*TopMiner, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get top miners failed because context failed")
	}

	space := userEconomySpace()
	params := map[string]interface{}{
		"limit":  limit,
		"offset": offset,
	}

	sql := fmt.Sprintf(`SELECT user_id, username, profile_picture_url, balance 
		FROM %[1]v INDEXED BY "pk_unnamed_%[1]v_1" 
		ORDER BY balance LIMIT :limit OFFSET :offset`, space)

	var res []*TopMiner
	if err := e.db.PrepareExecuteTyped(sql, params, &res); err != nil {
		return nil, errors.Wrapf(err, "failed to get %q record with limit:%v and offset:%v", space, limit, offset)
	}

	return res, nil
}
