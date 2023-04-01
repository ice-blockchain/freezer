// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"strings"

	"github.com/pkg/errors"
)

func (r *repository) initializeWorker(ctx context.Context, table, userID string, workerIndex int16) (err error) {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := fmt.Sprintf(`INSERT INTO %v(worker_index,user_id) VALUES ($1,$2)
						ON CONFLICT (worker_index,user_id) 
								DO NOTHING`, table)
	_, err = storagev2.Exec(ctx, r.dbV2, sql, workerIndex, userID)

	return errors.Wrapf(err, "failed to %v, for userID:%v", sql, userID)
}

func (r *repository) updateWorkerFields( //nolint:funlen // .
	ctx context.Context, workerIndex int16, table string, updateKV map[string]any, userIDs ...string,
) (err error) {
	if ctx.Err() != nil || len(userIDs) == 0 {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	ix := 0
	fields := make([]string, 0, len(updateKV))
	args := append(make([]any, 0, 1+1+len(updateKV)), workerIndex, userIDs)
	for key, value := range updateKV {
		fields = append(fields, fmt.Sprintf("%[1]v = $%[2]v", key, ix+1+1+1))
		args = append(args, value)
		ix++
	}
	sql := fmt.Sprintf(`UPDATE %[1]v
					    SET %[2]v
					    WHERE worker_index = $1 AND user_id = ANY($2)`, table, strings.Join(fields, ","))
	if _, uErr := storagev2.Exec(ctx, r.dbV2, sql, args...); uErr != nil {
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
	workerIndex = int16(uint64(resp.HashCode) % uint64(r.cfg.WorkerCount))

	return workerIndex, resp.HashCode, nil
}
