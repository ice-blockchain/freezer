// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage"
)

func (r *repository) initializeWorker(ctx context.Context, table, userID string, workerIndex uint64) (err error) {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`INSERT INTO %v%v(user_id) VALUES (:user_id)`, table, workerIndex)
	params := make(map[string]any, 1)
	params["user_id"] = userID
	if err = storage.CheckSQLDMLErr(r.db.PrepareExecute(sql, params)); err != nil && errors.Is(err, storage.ErrDuplicate) {
		return nil
	}

	return errors.Wrapf(err, "failed to %v, for userID:%v", sql, userID)
}

func (r *repository) updateWorkerFields(
	ctx context.Context, workerIndex uint64, table string, updateKV map[string]any, userIDs ...string,
) (err error) {
	if ctx.Err() != nil || len(userIDs) == 0 {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	values := make([]string, 0, len(userIDs))
	fields := make([]string, 0, len(updateKV))
	params := make(map[string]any, len(userIDs)+len(updateKV))
	for key, value := range updateKV {
		if value == nil {
			fields = append(fields, fmt.Sprintf("%[1]v = null", key))
		} else {
			params[key] = value
			fields = append(fields, fmt.Sprintf("%[1]v = :%[1]v", key))
		}
	}
	for i := range userIDs {
		params[fmt.Sprintf("user_id%v", i)] = userIDs[i]
		values = append(values, fmt.Sprintf(":user_id%v", i))
	}
	sql := fmt.Sprintf(`UPDATE %[1]v%[2]v
					    SET %[3]v
					    WHERE user_id in (%[4]v)`, table, workerIndex, strings.Join(fields, ","), strings.Join(values, ","))
	if _, uErr := storage.CheckSQLDMLResponse(r.db.PrepareExecute(sql, params)); uErr != nil {
		return errors.Wrapf(uErr, "failed to UPDATE %v%v params :%#v, for userIDs:%#v", table, workerIndex, params, userIDs)
	}

	return nil
}

func (r *repository) getWorkerIndex(ctx context.Context, userID string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, errors.Wrap(ctx.Err(), "context failed")
	}
	sql := `SELECT hash_code % :workers FROM users where user_id = :user_id`
	params := make(map[string]any, 1+1)
	params["workers"] = r.cfg.WorkerCount
	params["user_id"] = userID
	resp := make([]*struct {
		_msgpack    struct{} `msgpack:",asArray"` //nolint:tagliatelle,revive,nosnakecase // To insert we need asArray
		WorkerIndex uint64
	}, 0, 1)
	if err := r.db.PrepareExecuteTyped(sql, params, &resp); err != nil {
		return 0, errors.Wrapf(err, "failed to get worker index for userID:%v", userID)
	}
	if len(resp) == 0 {
		return 0, ErrNotFound
	}

	return resp[0].WorkerIndex, nil
}
