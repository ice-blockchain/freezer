// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync/atomic"
	stdlibtime "time"

	"github.com/cenkalti/backoff/v4"
	"github.com/goccy/go-json"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

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
	repo := &repository{
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
	go repo.startDisableAdvancedTeamCfgSyncer(ctx)

	return repo
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
	//nolint:contextcheck // It's intended. Cuz we want to close everything gracefully.
	mbConsumer := messagebroker.MustConnectAndStartConsuming(context.Background(), cancel, applicationYamlKey,
		&usersTableSource{processor: prc},
		&miningSessionsTableSource{processor: prc},
		&completedTasksSource{processor: prc},
		&viewedNewsSource{processor: prc},
		&deviceMetadataTableSource{processor: prc},
	)
	prc.shutdown = closeAll(mbConsumer, prc.mb, prc.db)

	go prc.startDisableAdvancedTeamCfgSyncer(ctx)
	go prc.startKYCConfigJSONSyncer(ctx)
	prc.mustInitAdoptions(ctx)
	prc.mustNotifyCurrentAdoption(ctx)
	prc.extraBonusStartDate = extrabonusnotifier.MustGetExtraBonusStartDate(ctx, prc.db)
	prc.extraBonusIndicesDistribution = extrabonusnotifier.MustGetExtraBonusIndicesDistribution(ctx, prc.db)
	prc.livenessLoadDistributionStartDate = mustGetLivenessLoadDistributionStartDate(ctx, prc.db)
	log.Info(fmt.Sprintf("configuration loaded[livenessLoadDistributionStartDate]: %#v", prc.livenessLoadDistributionStartDate))
	log.Info(fmt.Sprintf("configuration loaded[FaceRecognitionDelay]: %#v", cfg.KYC.FaceRecognitionDelay))
	log.Info(fmt.Sprintf("configuration loaded[LivenessDelay]: %#v", cfg.KYC.LivenessDelay))
	log.Info(fmt.Sprintf("configuration loaded[AdoptionMilestoneSwitch]: %#v", cfg.AdoptionMilestoneSwitch))
	log.Info(fmt.Sprintf("configuration loaded[ExtraBonuses]: %#v", cfg.ExtraBonuses))
	log.Info(fmt.Sprintf("configuration loaded[RollbackNegativeMining]: %#v", cfg.RollbackNegativeMining))
	log.Info(fmt.Sprintf("configuration loaded[MiningSessionDuration]: %#v", cfg.MiningSessionDuration))
	log.Info(fmt.Sprintf("configuration loaded[GlobalAggregationInterval]: %#v", cfg.GlobalAggregationInterval))

	return prc
}

func (r *repository) Close() error {
	return errors.Wrap(r.shutdown(), "closing repository failed")
}

func (r *repository) CheckHealth(ctx context.Context) error {
	return multierror.Append( //nolint:wrapcheck // Not needed.
		errors.Wrap(r.checkDBHealth(ctx), "db ping failed"),
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
	if err := p.checkDBHealth(ctx); err != nil {
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

func (r *repository) checkDBHealth(ctx context.Context) error {
	if resp := r.db.Ping(ctx); resp.Err() != nil || resp.Val() != "PONG" {
		if resp.Err() == nil {
			resp.SetErr(errors.Errorf("response `%v` is not `PONG`", resp.Val()))
		}

		return errors.Wrap(resp.Err(), "[health-check] failed to ping DB")
	}
	if !r.db.IsRW(ctx) {
		return errors.New("db is not writeable")
	}

	return nil
}

func (r *repository) startDisableAdvancedTeamCfgSyncer(ctx context.Context) {
	ticker := stdlibtime.NewTicker(5 * stdlibtime.Minute) //nolint:gosec,gomnd // Not an  issue.
	defer ticker.Stop()
	r.cfg.disableAdvancedTeam = new(atomic.Pointer[[]string])
	log.Panic(errors.Wrap(r.syncDisableAdvancedTeamCfg(ctx), "failed to syncDisableAdvancedTeamCfg"))

	for {
		select {
		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
			log.Error(errors.Wrap(r.syncDisableAdvancedTeamCfg(reqCtx), "failed to syncDisableAdvancedTeamCfg"))
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

func (r *repository) syncDisableAdvancedTeamCfg(ctx context.Context) error {
	result, err := r.db.Get(ctx, "disable_advanced_team_cfg").Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return errors.Wrap(err, "could not get `disable_advanced_team_cfg`")
	}
	var (
		oldCfg []string
		newCfg = strings.Split(strings.ReplaceAll(strings.ToLower(result), " ", ""), ",")
	)
	sort.SliceStable(newCfg, func(ii, jj int) bool { return newCfg[ii] < newCfg[jj] })
	if old := r.cfg.disableAdvancedTeam.Swap(&newCfg); old != nil {
		oldCfg = *old
	}
	if strings.Join(oldCfg, "") != strings.Join(newCfg, "") {
		log.Info(fmt.Sprintf("`disable_advanced_team_cfg` changed from: %#v, to: %#v", oldCfg, newCfg))
	}

	return nil
}

func (r *repository) isAdvancedTeamEnabled(device string) bool {
	if device == "" {
		return true
	}
	var disableAdvancedTeamFor []string
	if cfgVal := r.cfg.disableAdvancedTeam.Load(); cfgVal != nil {
		disableAdvancedTeamFor = *cfgVal
	}
	if len(disableAdvancedTeamFor) == 0 {
		return true
	}
	for _, disabled := range disableAdvancedTeamFor {
		if strings.EqualFold(device, disabled) {
			return false
		}
	}

	return true
}

func (r *repository) isAdvancedTeamDisabled(device string) bool {
	return !r.isAdvancedTeamEnabled(device)
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

func ContextWithAuthorization(ctx context.Context, authorization string) context.Context {
	if authorization == "" {
		return ctx
	}

	return context.WithValue(ctx, authorizationCtxValueKey, authorization) //nolint:revive,staticcheck // Not an issue.
}

func ContextWithXAccountMetadata(ctx context.Context, xAccountMetadata string) context.Context {
	if xAccountMetadata == "" {
		return ctx
	}

	return context.WithValue(ctx, xAccountMetadataCtxValueKey, xAccountMetadata) //nolint:revive,staticcheck // Not an issue.
}

func authorization(ctx context.Context) (authorization string) {
	authorization, _ = ctx.Value(authorizationCtxValueKey).(string) //nolint:errcheck // Not needed.

	return
}

func xAccountMetadata(ctx context.Context) (xAccountMetadata string) {
	xAccountMetadata, _ = ctx.Value(xAccountMetadataCtxValueKey).(string) //nolint:errcheck // Not needed.

	return
}

func ContextWithClientType(ctx context.Context, clientType string) context.Context {
	if clientType == "" {
		return ctx
	}

	return context.WithValue(ctx, clientTypeCtxValueKey, clientType) //nolint:revive,staticcheck // Not an issue.
}

func isWebClientType(ctx context.Context) bool {
	clientType, _ := ctx.Value(clientTypeCtxValueKey).(string) //nolint:errcheck // Not needed.

	return strings.EqualFold(strings.TrimSpace(clientType), "web")
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
