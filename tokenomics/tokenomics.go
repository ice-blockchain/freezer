// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	dwh "github.com/ice-blockchain/freezer/bookkeeper/storage"
	extrabonusnotifier "github.com/ice-blockchain/freezer/extra-bonus-notifier"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/multimedia/picture"
	"github.com/ice-blockchain/wintr/time"
)

func New(ctx context.Context, _ context.CancelFunc) Repository {
	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, applicationYamlKey)
	dwhClient := dwh.MustConnect(ctx, applicationYamlKey)

	return &repository{
		cfg:                           &cfg,
		extraBonusStartDate:           extrabonusnotifier.MustGetExtraBonusStartDate(ctx, db),
		extraBonusIndicesDistribution: extrabonusnotifier.MustGetExtraBonusIndicesDistribution(ctx, db),
		shutdown: func() error {
			return multierror.Append(db.Close(), dwhClient.Close()).ErrorOrNil()
		},
		db:            db,
		dwh:           dwhClient,
		pictureClient: picture.New(applicationYamlKey),
	}
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor {
	var cfg Config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	prc := &processor{repository: &repository{
		cfg:           &cfg,
		db:            storage.MustConnect(context.Background(), applicationYamlKey),
		mb:            messagebroker.MustConnect(context.Background(), applicationYamlKey),
		pictureClient: picture.New(applicationYamlKey),
	}}

	prc.mustInitAdoptions(ctx)
	prc.mustInitAdoptionSwitchTime(ctx)
	prc.extraBonusStartDate = extrabonusnotifier.MustGetExtraBonusStartDate(ctx, prc.db)
	prc.extraBonusIndicesDistribution = extrabonusnotifier.MustGetExtraBonusIndicesDistribution(ctx, prc.db)

	//nolint:contextcheck // It's intended. Cuz we want to close everything gracefully.
	mbConsumer := messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey,
		&usersTableSource{processor: prc},
		&miningSessionsTableSource{processor: prc},
		&completedTasksSource{processor: prc},
		&viewedNewsSource{processor: prc},
		&deviceMetadataTableSource{processor: prc},
	)
	prc.shutdown = closeAll(mbConsumer, prc.mb, prc.db)
	prc.mustNotifyCurrentAdoption(ctx)

	return prc
}

func (r *repository) Close() error {
	return errors.Wrap(r.shutdown(), "closing repository failed")
}

func (r *repository) CheckHealth(ctx context.Context) error {
	return multierror.Append( //nolint:wrapcheck // Not needed.
		errors.Wrap(r.pingDB(ctx), "db ping failed"),
		errors.Wrap(r.dwh.Ping(ctx), "dwh ping failed"),
	).ErrorOrNil()
}

func closeAll(mbConsumer, mbProducer messagebroker.Client, db storage.DB, otherClosers ...func() error) func() error {
	return func() error {
		err1 := errors.Wrap(mbConsumer.Close(), "closing mbConsumer connection failed")
		err2 := errors.Wrap(db.Close(), "closing db connection failed")
		err3 := errors.Wrap(mbProducer.Close(), "closing message broker producer connection failed")
		errs := make([]error, 0, 1+1+1+len(otherClosers))
		errs = append(errs, err1, err2, err3)
		for _, closeOther := range otherClosers {
			if err := closeOther(); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Wrap(multierror.Append(nil, errs...).ErrorOrNil(), "failed to close resources")
	}
}

func (p *processor) Close() error {
	return errors.Wrap(p.repository.Close(), "closing repository failed")
}

func (p *processor) CheckHealth(ctx context.Context) error {
	if err := p.pingDB(ctx); err != nil {
		return err
	}
	type ts struct {
		TS *time.Time `json:"ts"`
	}
	now := ts{TS: time.Now()}
	bytes, err := json.MarshalContext(ctx, now)
	if err != nil {
		return errors.Wrapf(err, "[health-check] failed to marshal %#v", now)
	}
	responder := make(chan error, 1)
	p.mb.SendMessage(ctx, &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     p.cfg.MessageBroker.Topics[0].Name,
		Topic:   p.cfg.MessageBroker.Topics[0].Name,
		Value:   bytes,
	}, responder)

	return errors.Wrapf(<-responder, "[health-check] failed to send health check message to broker")
}

func (r *repository) pingDB(ctx context.Context) error {
	if resp := r.db.Ping(ctx); resp.Err() != nil || resp.Val() != "PONG" {
		if resp.Err() == nil {
			resp.SetErr(errors.Errorf("response `%v` is not `PONG`", resp.Val()))
		}

		return errors.Wrap(resp.Err(), "[health-check] failed to ping DB")
	}

	return nil
}

func retry(ctx context.Context, op func() error) error {
	//nolint:wrapcheck // No need, its just a proxy.
	return backoff.RetryNotify(
		op,
		//nolint:gomnd // Because those are static configs.
		backoff.WithContext(&backoff.ExponentialBackOff{
			InitialInterval:     100 * stdlibtime.Millisecond,
			RandomizationFactor: 0.5,
			Multiplier:          2.5,
			MaxInterval:         stdlibtime.Second,
			MaxElapsedTime:      25 * stdlibtime.Second,
			Stop:                backoff.Stop,
			Clock:               backoff.SystemClock,
		}, ctx),
		func(e error, next stdlibtime.Duration) {
			log.Error(errors.Wrapf(e, "call failed. retrying in %v... ", next))
		})
}

func ContextWithHashCode(ctx context.Context, hashCode uint64) context.Context {
	if hashCode == 0 {
		return ctx
	}

	return context.WithValue(ctx, userHashCodeCtxValueKey, hashCode) //nolint:revive,staticcheck // Not an issue.
}

func requestingUserID(ctx context.Context) (requestingUserID string) {
	requestingUserID, _ = ctx.Value(requestingUserIDCtxValueKey).(string) //nolint:errcheck // Not needed.

	return
}

func (c *Config) totalActiveUsersAggregationIntervalDateFormat() string {
	const hoursInADay = 24
	switch c.AdoptionMilestoneSwitch.Duration { //nolint:exhaustive // We don't care about the others.
	case stdlibtime.Minute:
		return minuteFormat
	case stdlibtime.Hour:
		return hourFormat
	case hoursInADay * stdlibtime.Hour:
		return dayFormat
	default:
		log.Panic(fmt.Sprintf("invalid interval: %v", c.AdoptionMilestoneSwitch.Duration))

		return ""
	}
}

func (c *Config) globalAggregationIntervalChildDateFormat() string {
	const hoursInADay = 24
	switch c.GlobalAggregationInterval.Child { //nolint:exhaustive // We don't care about the others.
	case stdlibtime.Minute:
		return minuteFormat
	case stdlibtime.Hour:
		return hourFormat
	case hoursInADay * stdlibtime.Hour:
		return dayFormat
	default:
		log.Panic(fmt.Sprintf("invalid interval: %v", c.GlobalAggregationInterval.Child))

		return ""
	}
}

func (c *Config) globalAggregationIntervalParentDateFormat() string {
	const hoursInADay = 24
	switch c.GlobalAggregationInterval.Parent { //nolint:exhaustive // We don't care about the others.
	case stdlibtime.Minute:
		return minuteFormat
	case stdlibtime.Hour:
		return hourFormat
	case hoursInADay * stdlibtime.Hour:
		return dayFormat
	default:
		log.Panic(fmt.Sprintf("invalid interval: %v", c.GlobalAggregationInterval.Parent))

		return ""
	}
}
