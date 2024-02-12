// SPDX-License-Identifier: ice License 1.0

package quiz

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/freezer/model"
	appcfg "github.com/ice-blockchain/wintr/config"
	"github.com/ice-blockchain/wintr/connectors/storage/v2"
)

func NewRepository(ctx context.Context, _ context.CancelFunc) Repository {
	var cfg config
	appcfg.MustLoadFromKey(applicationYamlKey, &cfg)
	r := &repository{
		cfg: &cfg,
		db:  storage.MustConnect(ctx, "", applicationYamlKey),
	}
	return r
}

func (r *repository) SyncQuizStatus(ctx context.Context, usersToSync map[int64]*Status, histories []*model.User) error {
	var params, placeholders, i = []any{r.cfg.MaxResetCount}, make([]string, 0, len(usersToSync)), 1
	for _, qs := range usersToSync {
		params = append(params, qs.UserID)
		placeholders = append(placeholders, fmt.Sprintf("$%v", i+1)) //nolint:gomnd // Not a magic number.
		i++
	}
	sql := fmt.Sprintf(`SELECT 
				u.id,
				(qr.user_id IS NOT NULL AND cardinality(qr.resets) > $1)                                                     AS kyc_quiz_disabled,
				(qs.user_id IS NOT NULL AND qs.ended_at is not null AND qs.ended_successfully = true)                        AS kyc_quiz_completed
				FROM users u
					LEFT JOIN quiz_resets qr 
						   ON qr.user_id = u.id
					LEFT JOIN quiz_sessions qs
						   ON qs.user_id = u.id
					WHERE u.id IN (%v)
		`, strings.Join(placeholders, ","))
	res, err := storage.Select[struct {
		UserID        string `db:"id"`
		QuizDisabled  bool   `db:"kyc_quiz_disabled"`
		QuizCompleted bool   `db:"kyc_quiz_completed"`
	}](ctx, r.db, sql, params...)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch quiz status")
	}
	for _, qs := range res {
		var id int64
		for idx := range histories {
			if qs.UserID == histories[idx].UserID {
				id = histories[idx].ID
			}
			histories[idx].KYCQuizCompleted = qs.QuizCompleted
			histories[idx].KYCQuizDisabled = qs.QuizDisabled
		}
		usersToSync[id] = &Status{
			UserIDField:           model.UserIDField{UserID: qs.UserID},
			KYCQuizDisabledField:  model.KYCQuizDisabledField{KYCQuizDisabled: qs.QuizDisabled},
			KYCQuizCompletedField: model.KYCQuizCompletedField{KYCQuizCompleted: qs.QuizCompleted},
		}
	}

	return nil
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) CheckHealth(ctx context.Context) error {
	return errors.Wrap(r.db.Ping(ctx), "[health-check] quiz: failed to ping DB")
}
