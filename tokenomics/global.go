// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func (r *repository) getGlobalUnsignedValue(ctx context.Context, key string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, errors.Wrap(ctx.Err(), "context failed")
	}
	val, err := storagev2.Get[users.GlobalUnsigned](ctx, r.dbV2, `SELECT * FROM global WHERE key = $1`, key)
	if err != nil {
		return 0, errors.Wrapf(err, "failed to get global value for key:%v ", key)
	}

	return val.Value, nil
}

func (s *globalTableSource) Process(ctx context.Context, msg *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline while processing message")
	}
	if len(msg.Value) == 0 {
		return nil
	}
	var val users.GlobalUnsigned
	if err := json.UnmarshalContext(ctx, msg.Value, &val); err != nil {
		return errors.Wrapf(err, "process: cannot unmarshall %v into %#v", string(msg.Value), &val)
	}
	if val.Key == "" {
		return nil
	}
	sql := `INSERT INTO global(key,value) VALUES($1,$2::bigint)
			ON CONFLICT (key)
					 DO UPDATE
						   SET value = EXCLUDED.value
					 WHERE value != EXCLUDED.value`
	_, err := storagev2.Exec(ctx, s.dbV2, sql, val.Key, val.Value)

	return errors.Wrapf(err, "failed to upsert global unsigned value:%#v", &val)
}
