// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"context"
	"fmt"
	"strings"
	stdlibtime "time"

	"github.com/pkg/errors"

	"github.com/ice-blockchain/wintr/connectors/storage/v2"
	"github.com/ice-blockchain/wintr/time"
)

type (
	TriggeringRecord struct {
		CreatedAt          *time.Time
		Username           string
		ReferredByUsername string
		UserID             string
		EarnerUserID       string
		EthAddress         string
		InternalID         int64
		Balance            float64
	}
)

func TriggerCoinDistribution(ctx context.Context, db storage.Execer, records []*TriggeringRecord) error {
	const columns = 9
	values := make([]string, 0, len(records))
	args := make([]any, 0, len(records)*columns)
	for ix, record := range records {
		values = append(values, generateValuesSQLParams(ix, columns))
		args = append(args,
			record.CreatedAt.Time,
			record.CreatedAt.Time,
			record.InternalID,
			int64(record.Balance*100),
			record.Username,
			record.ReferredByUsername,
			record.UserID,
			record.EarnerUserID,
			record.EthAddress)
	}
	sql := fmt.Sprintf(`INSERT INTO coin_distributions_by_earner(created_at,day,internal_id,balance,username,referred_by_username,user_id,earner_user_id,eth_address) 
																 VALUES %v
						ON CONFLICT (day, user_id, earner_user_id) DO UPDATE
							SET balance = EXCLUDED.balance`, strings.Join(values, ",\n"))
	_, err := storage.Exec(ctx, db, sql, args...)

	return errors.Wrapf(err, "failed to insert into coin_distributions_by_earner %#v", records)
}

func generateValuesSQLParams(index, columns int) string {
	params := make([]string, 0, columns)
	for ii := 1; ii <= columns; ii++ {
		params = append(params, fmt.Sprintf("$%v", index*columns+ii))
	}

	return fmt.Sprintf("(%v)", strings.Join(params, ","))
}

func X(ctx context.Context, db storage.Execer) error {
	sql := `INSERT INTO global(key,value) VALUES ('latest_processing_date',$1)
				ON CONFLICT (key) DO UPDATE
					SET value = EXCLUDED.value
				WHERE value != EXCLUDED.value`
	_, err := storage.Exec(ctx, db, sql, time.Now().Format(stdlibtime.DateOnly))

	return errors.Wrap(err, "failed to XXX ")
}
