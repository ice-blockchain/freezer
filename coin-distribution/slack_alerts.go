// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func (r *repository) sendNewCoinDistributionsAvailableForReviewSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":eyes:`%v` <%v|new coin distributions are available for review> :eyes:", r.cfg.Environment, r.cfg.ReviewURL)

	return errors.Wrap(sendSlackMessage(ctx, text, r.cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendNewCoinDistributionsAvailableForReviewSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":eyes:`%v` <%v|new coin distributions are available for review> :eyes:", cfg.Environment, cfg.ReviewURL)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func (r *repository) sendCurrentCoinDistributionsAvailableForReviewAreApprovedSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":white_check_mark:`%v` current pending coin distributions are approved and are going to be processed as soon as the coin-distributer comes online :white_check_mark:", r.cfg.Environment) //nolint:lll // .

	return errors.Wrap(sendSlackMessage(ctx, text, r.cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func (r *repository) sendCurrentCoinDistributionsAvailableForReviewAreDeniedSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":no_entry:`%v` current pending coin distributions are denied and will not be processed :no_entry:", r.cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, r.cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendCurrentCoinDistributionsFinishedBeingSentToEthereumSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":rocket:`%v` all pending coin distributions have been processed successfully :rocket:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendCoinDistributionsProcessingStoppedDueToUnrecoverableFailureSlackMessage(ctx context.Context, reason string) error {
	text := fmt.Sprintf(":bangbang:`%v` coin distribution processing stopped due to failure :bangbang:\n:rotating_light: reason: `%v` :rotating_light:", cfg.Environment, reason) //nolint:lll // .

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

//nolint:funlen // .
func sendSlackMessage(ctx context.Context, text, alertSlackWebhook string) error {
	message := struct {
		Text string `json:"text,omitempty"`
	}{
		Text: text,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return errors.Wrapf(err, "failed to Marshal slack message:%#v", message)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, alertSlackWebhook, bytes.NewBuffer(data))
	if err != nil {
		return errors.Wrap(err, "newRequestWithContext failed")
	}

	resp, err := new(http.Client).Do(req)
	if err != nil {
		return errors.Wrap(err, "slack webhook request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected statusCode:%v", resp.StatusCode)
	}

	return errors.Wrap(resp.Body.Close(), "failed to close body")
}
