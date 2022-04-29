// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"github.com/framey-io/go-tarantool"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)

	return &repository{
		close:          closeDB(db),
		ReadRepository: &economy{db: db},
	}
}

func closeAll(db tarantool.Connector, mb messagebroker.Client) func() error {
	return func() error {
		err1 := errors.Wrap(db.Close(), "closing db connection failed")
		err2 := errors.Wrap(mb.Close(), "closing message broker connection failed")
		if err1 != nil && err2 != nil {
			return multierror.Append(err1, err2)
		}
		var err error
		if err1 != nil {
			err = err1
		}
		if err2 != nil {
			err = err2
		}

		return errors.Wrapf(err, "failed to close all resources")
	}
}

func closeDB(db tarantool.Connector) func() error {
	return func() error {
		return errors.Wrap(db.Close(), "closing db connection failed")
	}
}

func (r *repository) Close() error {
	return errors.Wrap(r.close(), "closing economy repository failed")
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)
	mb := messagebroker.MustConnect(ctx, applicationYamlKey)

	return &processor{
		close:           closeAll(db, mb),
		ReadRepository:  &economy{db: db},
		WriteRepository: &economy{db: db, mb: mb},
	}
}

func (p *processor) Close() error {
	return errors.Wrap(p.close(), "closing economy processor failed")
}

func (p *processor) CheckHealth(ctx context.Context) error {
	//nolint:nolintlint // TODO implement me.
	return nil
}
