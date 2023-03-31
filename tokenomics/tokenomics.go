// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"math/rand"
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
	storagev2 "github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/multimedia/picture"
	"github.com/ice-blockchain/wintr/time"
)

func New(ctx context.Context, cancel context.CancelFunc) Repository {
	var cfg config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)

	db := storage.MustConnect(ctx, cancel, getDDL(ddl, &cfg), applicationYamlKey)
	dbV2 := storagev2.MustConnect(ctx, "", applicationYamlKey)

	return &repository{
		cfg: &cfg,
		shutdown: func() error {
			return multierror.Append(db.Close(), dbV2.Close())
		},
		db:            db,
		dbV2:          dbV2,
		pictureClient: picture.New(applicationYamlKey),
	}
}

func StartProcessor(ctx context.Context, cancel context.CancelFunc) Processor { //nolint:funlen // A lot of startup & shutdown ceremony.
	var cfg config
	appCfg.MustLoadFromKey(applicationYamlKey, &cfg)
	var mbConsumer messagebroker.Client
	prc := &processor{repository: &repository{
		cfg: &cfg,
		db: storage.MustConnect(context.Background(), func() { //nolint:contextcheck // It's intended. Cuz we want to close everything gracefully.
			if mbConsumer != nil {
				log.Error(errors.Wrap(mbConsumer.Close(), "failed to close mbConsumer due to db premature cancellation"))
			}
			cancel()
		}, getDDL(ddl, &cfg), applicationYamlKey),
		dbV2:          storagev2.MustConnect(ctx, getDDL(ddlV2, &cfg), applicationYamlKey),
		mb:            messagebroker.MustConnect(ctx, applicationYamlKey),
		pictureClient: picture.New(applicationYamlKey),
	}}
	//nolint:contextcheck // It's intended. Cuz we want to close everything gracefully.
	mbConsumer = messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey,
		&usersTableSource{processor: prc},
		&globalTableSource{processor: prc},
		&miningSessionsTableSource{processor: prc},
		&addBalanceCommandsSource{processor: prc},
		&viewedNewsSource{processor: prc},
		&deviceMetadataTableSource{processor: prc},
	)
	prc.shutdown = closeAll(mbConsumer, prc.mb, prc.db, prc.dbV2)

	prc.initializeExtraBonusWorkers()
	prc.mustNotifyCurrentAdoption(ctx)
	go prc.startStreams(ctx)
	go prc.startCleaners(ctx)

	return prc
}

func getDDL(ddl string, cfg *config) string {
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

func closeAll(mbConsumer, mbProducer messagebroker.Client, db tarantool.Connector, dbV2 *storagev2.DB, otherClosers ...func() error) func() error {
	return func() error {
		err1 := errors.Wrap(mbConsumer.Close(), "closing mbConsumer connection failed")
		err2 := errors.Wrap(db.Close(), "closing db connection failed")
		err3 := errors.Wrap(dbV2.Close(), "closing dbV2 connection failed")
		err4 := errors.Wrap(mbProducer.Close(), "closing message broker producer connection failed")
		errs := make([]error, 0, 1+1+1+1+len(otherClosers))
		errs = append(errs, err1, err2, err3, err4)
		for _, closeOther := range otherClosers {
			if err := closeOther(); err != nil {
				errs = append(errs, err)
			}
		}

		return errors.Wrap(multierror.Append(nil, errs...).ErrorOrNil(), "failed to close resources")
	}
}

func (p *processor) Close() error {
	if p.cancelStreams != nil {
		p.cancelStreams()
	}
	if p.streamsDoneWg != nil {
		p.streamsDoneWg.Wait()
	}

	return errors.Wrap(p.repository.Close(), "closing repository failed")
}

func (p *processor) startStreams(ctx context.Context) { //nolint:funlen // .
	if ctx.Err() != nil {
		return
	}
	log.Info("trying to start streams")
	const key = "streams-processing-exclusive-lock"
	tuple := &users.GlobalUnsigned{Key: key, Value: uint64(time.Now().UnixNano())}
	if err := storage.CheckNoSQLDMLErr(p.db.InsertTyped("GLOBAL", tuple, &[]*users.GlobalUnsigned{})); err != nil {
		log.Error(errors.Wrapf(err, "failed to start streams, because failed to insert into global: %#v", tuple))
		const waitDuration = 5 * stdlibtime.Second
		stdlibtime.Sleep(waitDuration)
		p.startStreams(ctx)

		return
	}
	log.Info("streams started")
	defer func() {
		log.Error(errors.Wrapf(storage.CheckNoSQLDMLErr(p.db.DeleteTyped("GLOBAL", "pk_unnamed_GLOBAL_1", tarantool.StringKey{S: key}, &[]*users.GlobalUnsigned{})), "failed to delete GLOBAL(%v)", key)) //nolint:lll // .
	}()
	streamsCtx, cancelStreams := context.WithCancel(ctx)
	p.streamsDoneWg = new(sync.WaitGroup)
	p.streamsDoneWg.Add(1 + 1 + 1 + 1)
	p.cancelStreams = cancelStreams
	go func() {
		defer p.streamsDoneWg.Done()
		(&balanceRecalculationTriggerStreamSource{processor: p}).start(streamsCtx)
	}()
	go func() {
		defer p.streamsDoneWg.Done()
		(&miningRatesRecalculationTriggerStreamSource{processor: p}).start(streamsCtx)
	}()
	go func() {
		defer p.streamsDoneWg.Done()
		(&blockchainBalanceSynchronizationTriggerStreamSource{processor: p}).start(streamsCtx)
	}()
	go func() {
		defer p.streamsDoneWg.Done()
		(&extraBonusProcessingTriggerStreamSource{processor: p}).start(streamsCtx)
	}()
	p.streamsDoneWg.Wait()
	log.Info("streams stopped")
}

func (p *processor) startCleaners(ctx context.Context) {
	ticker := stdlibtime.NewTicker(stdlibtime.Duration(10+rand.Intn(30)) * stdlibtime.Second) //nolint:gosec,gomnd // Not an  issue.
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			const deadline = 30 * stdlibtime.Second
			reqCtx, cancel := context.WithTimeout(ctx, deadline)
			log.Error(errors.Wrap(p.deleteOldProcessedMiningSessions(reqCtx), "failed to deleteOldProcessedMiningSessions"))
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (p *processor) deleteOldProcessedMiningSessions(ctx context.Context) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	sql := `DELETE FROM processed_mining_sessions WHERE session_number < :session_number`
	params := make(map[string]any, 1)
	params["session_number"] = p.sessionNumber(time.New(time.Now().Add(-24 * stdlibtime.Hour)))
	if _, err := storage.CheckSQLDMLResponse(p.db.PrepareExecute(sql, params)); err != nil {
		return errors.Wrap(err, "failed to delete old data from processed_mining_sessions")
	}

	return nil
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

func executeBatchConcurrently[ARG any](ctx context.Context, process func(context.Context, ARG) error, args []ARG) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "unexpected deadline")
	}
	if len(args) == 0 {
		return nil
	}
	wg := new(sync.WaitGroup)
	wg.Add(len(args))
	errChan := make(chan error, len(args))
	for i := range args {
		go func(ix int) {
			defer wg.Done()
			errChan <- errors.Wrapf(process(ctx, args[ix]), "failed to process:%#v", args[ix])
		}(i)
	}
	wg.Wait()
	close(errChan)
	errs := make([]error, 0, len(args))
	for err := range errChan {
		errs = append(errs, err)
	}

	return errors.Wrap(multierror.Append(nil, errs...).ErrorOrNil(), "at least one arg processing failed")
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
