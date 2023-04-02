// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func (r *repository) initializeWorker(ctx context.Context, table, userID string, hashCode uint64) (err error) {
	sql := fmt.Sprintf(`INSERT INTO %v(worker_index,user_id,hash_code) VALUES ($1,$2,$3) ON CONFLICT (worker_index,user_id) DO NOTHING`, table)
	_, err = storage.Exec(ctx, r.db, sql, int16(hashCode%uint64(r.cfg.WorkerCount)), userID, int64(hashCode))

	return errors.Wrapf(err, "failed to %v, for userID:%v,hashCode:%v", sql, userID, hashCode)
}

func (r *repository) updateWorkerFields(
	ctx context.Context, workerIndex int16, table string, updateKV map[string]any, userIDs ...string,
) (err error) {
	if len(userIDs) == 0 {
		return nil
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
					    WHERE worker_index = $1 
					      AND user_id = ANY($2)`, table, strings.Join(fields, ","))
	_, err = storage.Exec(ctx, r.db, sql, args...)

	return errors.Wrapf(err, "failed to UPDATE %v%v updateKV :%#v, for userIDs:%#v", table, workerIndex, updateKV, userIDs)
}

func (r *repository) getWorker(ctx context.Context, userID string) (workerIndex int16, hashCode int64, err error) {
	resp, err := storage.Get[struct{ HashCode int64 }](ctx, r.db, `SELECT hash_code FROM users WHERE user_id = $1`, userID)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "failed to get worker for userID:%v", userID)
	}

	return int16(uint64(resp.HashCode) % uint64(r.cfg.WorkerCount)), resp.HashCode, nil
}
