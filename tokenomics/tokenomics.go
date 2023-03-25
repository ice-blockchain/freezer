// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"strings"
	"sync"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/go-tarantool-client"
	appCfg "github.com/ice-blockchain/wintr/config"
	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/connectors/storage"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/multimedia/picture"
	"github.com/ice-blockchain/wintr/time"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	var cfg config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, cancel, getDDL(&cfg), applicationYamlKey)

	return &repository{
		cfg:           &cfg,
		shutdown:      db.Close,
		db:            db,
		pictureClient: picture.New(applicationYamlKey),
	}
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor { //nolint:funlen // A lot of startup & shutdown ceremony.
	var cfg config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	mbConsumers := make([]messagebroker.Client, 0, 1+1+1+1+1)
	prc := &processor{repository: &repository{
		cfg: &cfg,
		db: storage.MustConnect(context.Background(), func() { //nolint:contextcheck // It's intended. Cuz we want to close everything gracefully.
			consumerStops := make([]func() error, 0, len(mbConsumers))
			for _, mbConsumer := range mbConsumers {
				consumerStops = append(consumerStops, mbConsumer.Close)
			}
			log.Error(errors.Wrap(executeConcurrently(consumerStops...), "failed to close mbConsumers due to db premature cancellation"))
			cancel()
		}, getDDL(&cfg), applicationYamlKey),
		mb:            messagebroker.MustConnect(ctx, applicationYamlKey),
		pictureClient: picture.New(applicationYamlKey),
	}}
	//nolint:contextcheck // It's intended. Cuz we want to close everything gracefully.
	mbConsumers = append(mbConsumers,
		messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey,
			&usersTableSource{processor: prc},
			&globalTableSource{processor: prc},
			&miningSessionsTableSource{processor: prc},
			&addBalanceCommandsSource{processor: prc},
			&viewedNewsSource{processor: prc},
			&deviceMetadataTableSource{processor: prc},
		),
		messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey+"1",
			&balanceRecalculationTriggerStreamSource{processor: prc},
		),
		messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey+"2",
			&miningRatesRecalculationTriggerStreamSource{processor: prc},
		),
		messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey+"3",
			&blockchainBalanceSynchronizationTriggerStreamSource{processor: prc},
		),
		messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey+"4",
			&extraBonusProcessingTriggerStreamSource{processor: prc},
		),
	)
	prc.shutdown = closeAll(mbConsumers, prc.mb, prc.db)

	prc.initializeExtraBonusWorkers()
	prc.mustNotifyCurrentAdoption(ctx)
	go prc.startBalanceRecalculationTriggerSeedingStream(ctx)
	go prc.startMiningRatesRecalculationTriggerSeedingStream(ctx)
	go prc.startBlockchainBalanceSynchronizationTriggerSeedingStream(ctx)
	go prc.startExtraBonusProcessingTriggerSeedingStream(ctx)

	return prc
}

func getDDL(cfg *config) string {
	extraBonusesValues := make([]string, 0, len(cfg.ExtraBonuses.FlatValues))
	for ix, value := range cfg.ExtraBonuses.FlatValues {
		extraBonusesValues = append(extraBonusesValues, fmt.Sprintf("(%v,%v)", ix, value))
	}
	now := time.Now()
	adoptionStart := now.Add(-24 * stdlibtime.Hour).UnixNano()
	dailyBonusStart := now.Add(-1 * users.NanosSinceMidnight(now)).UnixNano()
	args := make([]any, 0, len(cfg.AdoptionMilestoneSwitch.ActiveUserMilestones)+1+1+1+1)
	args = append(args, adoptionStart, cfg.WorkerCount-1, strings.Join(extraBonusesValues, ","), dailyBonusStart)
	for ix := range cfg.AdoptionMilestoneSwitch.ActiveUserMilestones {
		args = append(args, cfg.AdoptionMilestoneSwitch.ActiveUserMilestones[ix])
	}

	return fmt.Sprintf(ddl, args...)
}

func (r *repository) Close() error {
	return errors.Wrap(r.shutdown(), "closing repository failed")
}

func closeAll(mbConsumers []messagebroker.Client, mbProducer messagebroker.Client, db tarantool.Connector, otherClosers ...func() error) func() error {
	return func() error {
		consumerStops := make([]func() error, 0, len(mbConsumers))
		for _, mbConsumer := range mbConsumers {
			consumerStops = append(consumerStops, mbConsumer.Close)
		}
		err1 := errors.Wrap(executeConcurrently(consumerStops...), "closing message broker consumers connection failed")
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

func (p *processor) CheckHealth(ctx context.Context) error {
	if _, err := p.db.Ping(); err != nil {
		return errors.Wrap(err, "[health-check] failed to ping DB")
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

func (r *repository) workerIndex(ctx context.Context) (workerIndex uint64) {
	userHashCode, _ := ctx.Value(userHashCodeCtxValueKey).(uint64) //nolint:errcheck // Not needed.

	return userHashCode % r.cfg.WorkerCount
}

func executeConcurrently(fs ...func() error) error {
	if len(fs) == 0 {
		return nil
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(fs))
	errChan := make(chan error, len(fs))
	for i := range fs {
		go func(ix int) {
			defer wg.Done()
			errChan <- errors.Wrapf(fs[ix](), "failed to run func with index [%v]", ix)
		}(i)
	}
	wg.Wait()
	close(errChan)
	errs := make([]error, 0, len(fs))
	for err := range errChan {
		errs = append(errs, err)
	}

	return errors.Wrap(multierror.Append(nil, errs...).ErrorOrNil(), "at least one execution failed")
}

func sendMessagesConcurrently[M any](ctx context.Context, sendMessage func(context.Context, M) error, messages []M) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	if len(messages) == 0 {
		return nil
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(messages))
	errChan := make(chan error, len(messages))
	for i := range messages {
		go func(ix int) {
			defer wg.Done()
			errChan <- errors.Wrapf(sendMessage(ctx, messages[ix]), "failed to sendMessage:%#v", messages[ix])
		}(i)
	}
	wg.Wait()
	close(errChan)
	errs := make([]error, 0, len(messages))
	for err := range errChan {
		errs = append(errs, err)
	}

	return errors.Wrap(multierror.Append(nil, errs...).ErrorOrNil(), "at least one message sends failed")
}

func (c *config) globalAggregationIntervalChildDateFormat() string {
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

func (c *config) globalAggregationIntervalParentDateFormat() string {
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

func (c *config) lastXMiningSessionsCollectingIntervalDateFormat() string {
	const hoursInADay = 24
	switch c.RollbackNegativeMining.LastXMiningSessionsCollectingInterval { //nolint:exhaustive // We don't care about the others.
	case stdlibtime.Minute:
		return minuteFormat
	case stdlibtime.Hour:
		return hourFormat
	case hoursInADay * stdlibtime.Hour:
		return dayFormat
	default:
		log.Panic(fmt.Sprintf("invalid interval: %v", c.RollbackNegativeMining.LastXMiningSessionsCollectingInterval))

		return ""
	}
}

func (r *repository) lastXMiningSessionsCollectingIntervalDateFormat(now *time.Time) string {
	return now.Format(r.cfg.lastXMiningSessionsCollectingIntervalDateFormat())
}
