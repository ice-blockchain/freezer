// SPDX-License-Identifier: ice License 1.0

package coindistribution

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	stdlibtime "time"

	"github.com/goccy/go-json"
	"github.com/pkg/errors"
)

func (r *repository) sendCurrentCoinDistributionsAvailableForReviewAreApprovedSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":white_check_mark:`%v` current pending coin distributions are approved and are going to be processed as soon as the coin-distributer comes online :white_check_mark:", r.cfg.Environment) //nolint:lll // .

	return errors.Wrap(sendSlackMessage(ctx, text, r.cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func (r *repository) sendCurrentCoinDistributionsAvailableForReviewAreDeniedSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":no_entry:`%v` current pending coin distributions are denied and will not be processed :no_entry:", r.cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, r.cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendNewCoinDistributionsAvailableForReviewSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":eyes:`%v` <%v|new coin distributions are available for review> :eyes:", cfg.Environment, cfg.ReviewURL)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func SendNewCoinDistributionCollectionCycleStartedSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":money_mouth_face:`%v` started to collect coins for ethereum distribution :money_mouth_face:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func SendNewCoinDistributionCollectionCycleEndedPrematurelySlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":recycle:`%v` collecting coins for ethereum distribution stopped prematurely :recycle:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendCoinDistributerIsNowOnlineSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":sun_with_face:`%v` coin distributer is now online :sun_with_face:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendCoinDistributerIsNowOfflineSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":sleeping:`%v` coin distributer is now offline :sleeping:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendCurrentCoinDistributionsFinishedBeingSentToEthereumSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":rocket:`%v` all pending coin distributions have been processed successfully :rocket:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendAllCurrentCoinDistributionsWereCommittedInEthereumSlackMessage(ctx context.Context) error {
	text := fmt.Sprintf(":tada:`%v` all coin distributions have been committed successfully in ethereum :tada:", cfg.Environment)

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

func sendCoinDistributionsProcessingStoppedDueToUnrecoverableFailureSlackMessage(ctx context.Context, reason string) error {
	text := fmt.Sprintf(":bangbang:`%v` coin distribution processing stopped due to failure :bangbang:\n:rotating_light: reason: `%v` :rotating_light:", cfg.Environment, reason) //nolint:lll // .

	return errors.Wrap(sendSlackMessage(ctx, text, cfg.AlertSlackWebhook), "failed to sendSlackMessage")
}

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

	const retries = 10
	var resp *http.Response
	for ix := 0; ix < retries; ix++ {
		if resp, err = new(http.Client).Do(req); err == nil && resp.StatusCode == http.StatusOK {
			break
		}
		stdlibtime.Sleep(stdlibtime.Second)
	}
	if err != nil {
		return errors.Wrap(err, "slack webhook request failed")
	}
	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("unexpected statusCode:%v", resp.StatusCode)
	}

	return errors.Wrap(resp.Body.Close(), "failed to close body")
}
