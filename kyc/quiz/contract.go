// SPDX-License-Identifier: ice License 1.0

package quiz

import (
	"context"
	"io"

	"github.com/ice-blockchain/freezer/model"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

type (
	Repository interface {
		io.Closer
		SyncQuizStatus(ctx context.Context, usersToFetchTheQuiz map[int64]*Status, histories []*model.User) error
		CheckHealth(ctx context.Context) error
	}

	Status struct {
		model.UserIDField
		model.KYCQuizDisabledField
		model.KYCQuizCompletedField
	}
)

const (
	applicationYamlKey = "kyc/quiz"
)

type (
	config struct {
		MaxResetCount int64 `yaml:"maxResetCount"`
	}

	repository struct {
		cfg *config
		db  *storage.DB
	}
)
