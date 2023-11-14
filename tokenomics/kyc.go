// SPDX-License-Identifier: ice License 1.0

package tokenomics

import (
	"context"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"

	"github.com/ice-blockchain/eskimo/users"
	"github.com/ice-blockchain/wintr/connectors/storage/v3"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/terror"
	"github.com/ice-blockchain/wintr/time"
)

func init() { //nolint:gochecknoinits // It's the only way to tweak the client.
	req.DefaultClient().SetJsonMarshal(json.Marshal)
	req.DefaultClient().SetJsonUnmarshal(json.Unmarshal)
	req.DefaultClient().GetClient().Timeout = requestDeadline
}

func (r *repository) startKYCConfigJSONSyncer(ctx context.Context) {
	ticker := stdlibtime.NewTicker(stdlibtime.Minute) //nolint:gosec,gomnd // Not an  issue.
	defer ticker.Stop()
	r.cfg.kycConfigJSON = new(atomic.Pointer[kycConfigJSON])
	log.Panic(errors.Wrap(r.syncKYCConfigJSON(ctx), "failed to syncKYCConfigJSON"))

	for {
		select {
		case <-ticker.C:
			reqCtx, cancel := context.WithTimeout(ctx, requestDeadline)
			log.Error(errors.Wrap(r.syncKYCConfigJSON(reqCtx), "failed to syncKYCConfigJSON"))
			cancel()
		case <-ctx.Done():
			return
		}
	}
}

//nolint:funlen,gomnd // .
func (r *repository) syncKYCConfigJSON(ctx context.Context) error {
	if resp, err := req.
		SetContext(ctx).
		SetRetryCount(25).
		SetRetryBackoffInterval(10*stdlibtime.Millisecond, 1*stdlibtime.Second).
		SetRetryHook(func(resp *req.Response, err error) {
			if err != nil {
				log.Error(errors.Wrap(err, "failed to fetch KYCConfigJSON, retrying..."))
			} else {
				log.Error(errors.Errorf("failed to fetch KYCConfigJSON with status code:%v, retrying...", resp.GetStatusCode()))
			}
		}).
		SetRetryCondition(func(resp *req.Response, err error) bool {
			return err != nil || resp.GetStatusCode() != http.StatusOK
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Cache-Control", "no-cache, no-store, must-revalidate").
		SetHeader("Pragma", "no-cache").
		SetHeader("Expires", "0").
		Get(r.cfg.KYC.ConfigJSONURL); err != nil {
		return errors.Wrapf(err, "failed to get fetch `%v`", r.cfg.KYC.ConfigJSONURL)
	} else if data, err2 := resp.ToBytes(); err2 != nil {
		return errors.Wrapf(err2, "failed to read body of `%v`", r.cfg.KYC.ConfigJSONURL)
	} else {
		var kycConfig kycConfigJSON
		if err = json.UnmarshalContext(ctx, data, &kycConfig); err != nil {
			return errors.Wrapf(err, "failed to unmarshal into %#v, data: %v", kycConfig, string(data))
		}
		if kycConfig == (kycConfigJSON{}) {
			if body := string(data); !strings.Contains(body, "face-auth") && !strings.Contains(body, "web-face-auth") {
				return errors.Errorf("there's something wrong with the KYCConfigJSON body: %v", body)
			}
		}
		r.cfg.kycConfigJSON.Swap(&kycConfig)

		return nil
	}
}

func (r *repository) validateKYC(ctx context.Context, state *getCurrentMiningSession, skipKYCStep *users.KYCStep) error { //nolint:funlen // .
	if skipKYCStep != nil {
		if *skipKYCStep == users.FacialRecognitionKYCStep || *skipKYCStep == users.LivenessDetectionKYCStep {
			return errors.Errorf("you can't skip kycStep:%v", *skipKYCStep)
		}
		if *skipKYCStep == users.NoneKYCStep { // TODO implement this properly. This is used for mocking atm.
			if rand.Intn(2) == 0 {
				return terror.New(ErrKYCRequired, map[string]any{
					"kycStep": []users.KYCStep{users.Social1KYCStep, users.QuizKYCStep, users.QuizKYCStep, users.Social2KYCStep}[rand.Intn(4)],
				})
			}
		}
	}
	if state.KYCStepBlocked > 0 {
		return terror.New(ErrMiningDisabled, map[string]any{
			"kycStepBlocked": state.KYCStepBlocked,
		})
	}
	switch state.KYCStepPassed {
	case users.NoneKYCStep:
		if !state.MiningSessionSoloLastStartedAt.IsNil() && r.isKYCEnabled(ctx, users.FacialRecognitionKYCStep) {
			return terror.New(ErrKYCRequired, map[string]any{
				"kycStep": users.FacialRecognitionKYCStep,
			})
		}
	case users.FacialRecognitionKYCStep:
		if r.isKYCEnabled(ctx, users.LivenessDetectionKYCStep) {
			return terror.New(ErrKYCRequired, map[string]any{
				"kycStep": users.LivenessDetectionKYCStep,
			})
		}
	default:
		isAfterDelay := time.Now().Sub(*(*state.KYCStepsLastUpdatedAt)[users.LivenessDetectionKYCStep-1].Time) >= r.cfg.KYC.LivenessDelay
		isReservedForToday := int64((time.Now().Sub(*r.livenessLoadDistributionStartDate.Time)%r.cfg.KYC.LivenessDelay)/r.cfg.MiningSessionDuration.Max) == state.ID%int64(r.cfg.KYC.LivenessDelay/r.cfg.MiningSessionDuration.Max) //nolint:lll // .
		if (isAfterDelay || isReservedForToday) && r.isKYCEnabled(ctx, users.LivenessDetectionKYCStep) {
			return terror.New(ErrKYCRequired, map[string]any{
				"kycStep": users.LivenessDetectionKYCStep,
			})
		}
	}

	return nil
}

func (r *repository) isKYCEnabled(ctx context.Context, _ users.KYCStep) bool {
	var (
		kycConfig = r.cfg.kycConfigJSON.Load()
		isWeb     = isWebClientType(ctx)
	)

	if isWeb && !kycConfig.WebFaceAuth.Enabled {
		return false
	}

	if !isWeb && !kycConfig.FaceAuth.Enabled {
		return false
	}

	return true
}

func mustGetLivenessLoadDistributionStartDate(ctx context.Context, db storage.DB) (livenessLoadDistributionStartDate *time.Time) {
	livenessLoadDistributionStartDateString, err := db.Get(ctx, "liveness_load_distribution_start_date").Result()
	if err != nil && errors.Is(err, redis.Nil) {
		err = nil
	}
	log.Panic(errors.Wrap(err, "failed to get liveness_load_distribution_start_date"))
	if livenessLoadDistributionStartDateString != "" {
		livenessLoadDistributionStartDate = new(time.Time)
		log.Panic(errors.Wrapf(livenessLoadDistributionStartDate.UnmarshalText([]byte(livenessLoadDistributionStartDateString)), "failed to parse liveness_load_distribution_start_date `%v`", livenessLoadDistributionStartDateString)) //nolint:lll // .
		livenessLoadDistributionStartDate = time.New(livenessLoadDistributionStartDate.UTC())

		return
	}
	livenessLoadDistributionStartDate = time.New(time.Now().Truncate(24 * stdlibtime.Hour))
	set, sErr := db.SetNX(ctx, "liveness_load_distribution_start_date", livenessLoadDistributionStartDate, 0).Result()
	log.Panic(errors.Wrap(sErr, "failed to set liveness_load_distribution_start_date"))
	if !set {
		return mustGetLivenessLoadDistributionStartDate(ctx, db)
	}

	return livenessLoadDistributionStartDate
}