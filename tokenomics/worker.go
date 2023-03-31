// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"strings"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/time"
)

func (r *repository) initializeWorker(ctx context.Context, table, userID string, workerIndex int16) (err error) {
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

func (r *repository) updateWorkerFields( //nolint:funlen // .
	ctx context.Context, workerIndex int16, table string, updateKV map[string]any, userIDs ...string,
) (err error) {
	if ctx.Err() != nil || len(userIDs) == 0 {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	fields := make([]string, 0, len(updateKV))
	for key, value := range updateKV {
		switch typedValue := value.(type) {
		case *time.Time:
			fields = append(fields, fmt.Sprintf("%[1]v = %[2]v", key, typedValue.UnixNano()))
		case string:
			fields = append(fields, fmt.Sprintf("%[1]v = '%[2]v'", key, typedValue))
		default:
			if typedValue == nil {
				fields = append(fields, fmt.Sprintf("%[1]v = null", key))
			} else {
				fields = append(fields, fmt.Sprintf("%[1]v = %[2]v", key, typedValue))
			}
		}
	}
	values := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		values = append(values, fmt.Sprintf("'%v'", userID))
	}
	sql := fmt.Sprintf(`UPDATE %[1]v%[2]v
					    SET %[3]v
					    WHERE user_id in (%[4]v)`, table, workerIndex, strings.Join(fields, ","), strings.Join(values, ","))
	if _, uErr := storage.CheckSQLDMLResponse(r.db.Execute(sql)); uErr != nil {
		return errors.Wrapf(uErr, "failed to UPDATE %v%v updateKV :%#v, for userIDs:%#v", table, workerIndex, updateKV, userIDs)
	}

	return nil
}

func (r *repository) getWorker(ctx context.Context, userID string) (workerIndex int16, hashCode int64, err error) {
	sql := `SELECT hash_code 
			FROM users 
			where user_id = $1`
	resp, err := storagev2.Get[struct{ HashCode int64 }](ctx, r.dbV2, sql, userID)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to get worker for userID:%v", userID)
	}
	workerIndex = int16(resp.HashCode % int64(r.cfg.WorkerCount))

	return workerIndex, resp.HashCode, nil
}
