// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"github.com/framey-io/go-tarantool"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/economy/internal/storages/users"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)

	return &repository{
		close:          closeDB(db),
		ReadRepository: &economy{db: db},
	}
}

func closeAll(db tarantool.Connector, mbProducer, mbConsumer messagebroker.Client) func() error {
	return func() error {
		err1 := errors.Wrap(mbConsumer.Close(), "closing message broker consumer connection failed")
		err2 := errors.Wrap(mbProducer.Close(), "closing message broker producer connection failed")
		err3 := errors.Wrap(db.Close(), "closing db connection failed")
		errs := make([]error, 0, 1+1+1)
		if err1 != nil {
			errs = append(errs, err1)
		}
		if err2 != nil {
			errs = append(errs, err2)
		}
		if err3 != nil {
			errs = append(errs, err3)
		}
		if len(errs) > 1 {
			return multierror.Append(nil, errs...)
		} else if len(errs) == 1 {
			return errors.Wrapf(errs[0], "failed to close all resources")
		}

		return nil
	}
}

func closeDB(db tarantool.Connector) func() error {
	return func() error {
		return errors.Wrap(db.Close(), "closing db connection failed")
	}
}

func (r *repository) Close() error {
	log.Info("closing economy repository...")

	return errors.Wrap(r.close(), "closing economy repository failed")
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)
	mbProducer := messagebroker.MustConnect(ctx, applicationYamlKey)

	mbProcessors := processors(context.Background(), db)
	mbConsumer := messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey, mbProcessors)

	return &processor{
		close:           closeAll(db, mbProducer, mbConsumer),
		ReadRepository:  &economy{db: db},
		WriteRepository: &economy{db: db, mbProducer: mbProducer},
	}
}

func processors(ctx context.Context, db tarantool.Connector) map[messagebroker.Topic]messagebroker.Processor {
	return map[messagebroker.Topic]messagebroker.Processor{
		cfg.MessageBroker.ConsumingTopics[0]: users.New(db),
	}
}

func (p *processor) Close() error {
	log.Info("closing economy processor...")

	return errors.Wrap(p.close(), "closing economy processor failed")
}

func (p *processor) CheckHealth(ctx context.Context) error {
	//nolint:nolintlint // TODO implement me.
	return nil
}
