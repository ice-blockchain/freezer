// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
)

func (r *repository) getGlobalUnsignedValue(ctx context.Context, key string) (uint64, error) {
	val, err := storage.Get[users.GlobalUnsigned](ctx, r.db, `SELECT * FROM global WHERE key = $1`, key)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get global value for key:%v ", key)
	}

	return val.Value, nil
}

//nolint:revive // Not an issue atm.
func (r *repository) insertGlobalUnsignedValue(ctx context.Context, val *users.GlobalUnsigned, replace bool) error {
	var sql string
	if replace {
		sql = `INSERT INTO global(key,value) VALUES($1,$2::bigint)
			   ON CONFLICT (key)
					    DO UPDATE
							  SET value = EXCLUDED.value
						WHERE global.value != EXCLUDED.value`
	} else {
		sql = `INSERT INTO global(key,value) VALUES($1,$2::bigint)`
	}
	_, err := storage.Exec(ctx, r.db, sql, val.Key, val.Value)

	return errors.Wrapf(err, "failed to insert[replace:%v] global val:%#v", replace, val)
}

func (s *globalTableSource) Process(ctx context.Context, msg *messagebroker.Message) error {
	if ctx.Err() != nil || len(msg.Value) == 0 {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	var val users.GlobalUnsigned
	if err := json.UnmarshalContext(ctx, msg.Value, &val); err != nil || val.Key == "" {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &val)
	}

	return errors.Wrapf(s.insertGlobalUnsignedValue(ctx, &val, true), "failed to insertGlobalUnsignedValue:%#v", &val)
}
