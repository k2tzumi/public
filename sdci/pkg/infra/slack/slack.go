// Package slack is a service to operate Slack integrations.
package slack // import "cirello.io/exp/sdci/pkg/infra/slack"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"cirello.io/errors"
)

// Sends a message to a given Slack webhook.
func Send(webhookURL string, msg string) error {
	var payload bytes.Buffer
	err := json.NewEncoder(&payload).Encode(struct {
		Text string `json:"text"`
	}{Text: msg})
	if err != nil {
		return errors.E(err, "cannot encode slack message")
	}
	response, err := http.Post(webhookURL, "application/json", &payload)
	if err != nil {
		return errors.E(err, "cannot send slack message")
	}
	if _, err := io.Copy(ioutil.Discard, response.Body); err != nil {
		return errors.E(err, "cannot drain response body")
	}
	return nil
}
