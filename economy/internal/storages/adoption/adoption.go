package adoption

import (
	"context"

	"github.com/framey-io/go-tarantool"
	"github.com/pkg/errors"

	messagebroker "github.com/ice-blockchain/wintr/connectors/message_broker"
)

func New(db tarantool.Connector) messagebroker.Processor {
	return &adoptionSource{r: newRepository(db).(*repository)}
}

func (a *adoptionSource) Process(ctx context.Context, _ *messagebroker.Message) error {
	if ctx.Err() != nil {
		return errors.Wrap(ctx.Err(), "context failed")
	}
	if err := a.r.updateActiveUsers(ctx); err != nil {
		return errors.Wrapf(err, "adoption/adoptionSource: failed to update active users")
	}
	if err := a.switchActiveAdoption(ctx); err != nil {
		return errors.Wrapf(err, "adoption/adoptionSource: failed to switch active adoption/mining rate")
	}

	return errors.Wrapf(a.r.updateTotalUsersHistory(ctx), "adoption/adoptionSource: failed to update total users history")
}

func (a *adoptionSource) switchActiveAdoption(ctx context.Context) error {
	//nolint:nolintlint,gocritic // TODO: IF the last 168 consecutive hours from adoption_history.hour_timestamp have ALL been >= ANY adoption.total_active_users,
	// then adoption.active of that entry becomes true and the previous active adoption entry becomes false.
	return nil
}
