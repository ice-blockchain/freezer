// SPDX-License-Identifier: ice License 1.0

package miner

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
	"github.com/ice-blockchain/wintr/log"
	"github.com/ice-blockchain/wintr/time"
)

func didANewDayOffJustStart(now *time.Time, usr *user) *DayOffStarted {
	miningSessionDuration := cfg.MiningSessionDuration.Max
	if usr == nil ||
		usr.MiningSessionSoloStartedAt.IsNil() ||
		usr.MiningSessionSoloEndedAt.IsNil() ||
		usr.MiningSessionSoloLastStartedAt.IsNil() ||
		usr.BalanceLastUpdatedAt.IsNil() ||
		usr.MiningSessionSoloEndedAt.Before(*now.Time) ||
		usr.MiningSessionSoloLastStartedAt.Add(miningSessionDuration).After(*now.Time) {
		return nil
	}
	startedAt := time.New(usr.MiningSessionSoloLastStartedAt.Add((now.Sub(*usr.MiningSessionSoloLastStartedAt.Time) / miningSessionDuration) * miningSessionDuration)) //nolint:lll // .
	if usr.BalanceLastUpdatedAt.After(*startedAt.Time) {
		return nil
	}

	return &DayOffStarted{
		StartedAt:                   startedAt,
		EndedAt:                     time.New(startedAt.Add(miningSessionDuration)),
		UserID:                      usr.UserID,
		ID:                          fmt.Sprint(startedAt.UnixNano() / miningSessionDuration.Nanoseconds()),
		RemainingFreeMiningSessions: uint64(usr.MiningSessionSoloEndedAt.Sub(*now.Time) / miningSessionDuration),
		MiningStreak:                uint64(now.Sub(*usr.MiningSessionSoloStartedAt.Time) / miningSessionDuration),
	}
}

func dayOffStartedMessage(ctx context.Context, event *DayOffStarted) *messagebroker.Message {
	valueBytes, err := json.MarshalContext(ctx, event)
	log.Panic(errors.Wrapf(err, "failed to marshal %#v", event))

	return &messagebroker.Message{
		Headers: map[string]string{"producer": "freezer"},
		Key:     event.UserID,
		Topic:   cfg.MessageBroker.ProducingTopics[0].Name,
		Value:   valueBytes,
	}
}
