// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func (e *economy) GetTopMiners(ctx context.Context, arg *GetTopMinersArg) ([]*TopMiner, error) {
	if ctx.Err() != nil {
		return nil, errors.Wrap(ctx.Err(), "get top miners failed because context failed")
	}

	sql := fmt.Sprintf(`
		SELECT  u.user_id,
				u.username, 
				u.profile_picture_url, 
				b.amount 
		FROM BALANCES b
				JOIN USER_ECONOMY u
					on lower(u.username) LIKE :keyword ESCAPE '\'
					AND b.user_id = u.user_id
					AND lower(b.type) = 'total'
		ORDER BY b.amount_w3 DESC,
				 b.amount_w2 DESC,
				 b.amount_w1 DESC,
				 b.amount_w0 DESC 
		LIMIT %v OFFSET :offset`, arg.Limit)

	var res []*TopMiner
	if err := e.db.PrepareExecuteTyped(sql, arg.params(), &res); err != nil {
		return nil, errors.Wrapf(err, "failed to select for top miners with arg:%#v", arg)
	}

	return res, nil
}

func (arg *GetTopMinersArg) params() map[string]interface{} {
	return map[string]interface{}{
		"offset":  arg.Offset,
		"keyword": fmt.Sprintf("%%%v%%", strings.ReplaceAll(strings.ToLower(arg.Keyword), "_", "\\_")),
	}
}
