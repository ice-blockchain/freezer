// SPDX-License-Identifier: BUSL-1.1

package economy

import "context"

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	// TODO implement me
	return &repository{}
}

func (r *repository) Close() error {
	// TODO implement me
	return nil
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	// TODO implement me
	return &processor{}
}

func (p *processor) Close() error {
	// TODO implement me
	return nil
}

func (p *processor) CheckHealth(ctx context.Context) error {
	// TODO implement me
	return nil
}
