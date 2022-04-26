// SPDX-License-Identifier: BUSL-1.1

package economy

import (
	"context"

	"github.com/pkg/errors"

	appCfg "github.com/ICE-Blockchain/wintr/config"
	"github.com/ICE-Blockchain/wintr/connectors/storage"
	"github.com/ICE-Blockchain/wintr/log"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	db := storage.MustConnect(ctx, cancel, ddl, applicationYamlKey)
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	return &repository{
		close: db.Close,
		db:    db,
	}
}

func (r *repository) Close() error {
	log.Info("closing economy repository...")

	return errors.Wrap(r.close(), "closing economy repository failed")
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	//nolint:nolintlint // TODO implement me
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
