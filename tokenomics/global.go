// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/go-tarantool-client"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
)

func (r *repository) getGlobalUnsignedValue(ctx context.Context, key string) (uint64, error) {
	if ctx.Err() != nil {
		return 0, errors.Wrap(ctx.Err(), "context failed")
	}
	var val users.GlobalUnsigned
	if err := r.db.GetTyped("GLOBAL", "pk_unnamed_GLOBAL_1", tarantool.StringKey{S: key}, &val); err != nil {
		return 0, errors.Wrapf(err, "failed to get global value for key:%v ", key)
	}
	if val.Key == "" {
		return 0, storage.ErrNotFound
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
	ops := make([]tarantool.Op, 0, 1)
	ops = append(ops, tarantool.Op{Op: "=", Field: 1, Arg: val.Value})

	return errors.Wrapf(s.db.UpsertTyped("GLOBAL", &val, ops, &[]*users.GlobalUnsigned{}),
		"failed to upsert global unsigned value:%#v", &val)
}
