// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"github.com/framey-io/go-tarantool"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	appCfg "github.com/ICE-Blockchain/wintr/config"
	messagebroker "github.com/ICE-Blockchain/wintr/connectors/message_broker"
	"github.com/ICE-Blockchain/wintr/connectors/storage"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)
	mb := messagebroker.MustConnect(ctx, applicationYamlKey)

	return &repository{
		close:           closeAll(db, mb),
		ReadRepository:  &userEconomyRepository{db, mb},
		WriteRepository: &userEconomyRepository{db, mb},
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

func (r *repository) Close() error {

	if err := r.close(); err != nil {
		return errors.Wrap(err, "unable close repository")
	}

	return nil
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	//nolint:nolintlint // TODO implement me.
	return nil
}

func (p *processor) Close() error {
	//nolint:nolintlint // TODO implement me.
	return nil
}

func (p *processor) CheckHealth(ctx context.Context) error {
	//nolint:nolintlint // TODO implement me.
	return nil
}
