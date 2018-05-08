package rabbitmq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"io/ioutil"
	"net/http"
	"regexp"
)

func pushCommandResponse(response *rabbitmq.CommandResponse) bool {
	callbackURL := response.Channel
	if len(common.GConfig.Mattermost.ReplacementURL) > 0 {
		r := regexp.MustCompile(`^https?:\/\/.+\/(.+)$`)
		callbackURL = r.ReplaceAllString(
			response.Channel,
			fmt.Sprintf("%s/$1", common.GConfig.Mattermost.ReplacementURL),
		)
	}

	common.Log.Infof("Received command response for user %s (callback %s)", response.User, response.Channel)
	pushResponse := fmt.Sprintf(`{"response_type": "ephemeral", "text": "%s"}`, response.Message)
	req, err := http.NewRequest("POST", callbackURL, bytes.NewBuffer([]byte(pushResponse)))
	if err != nil {
		common.Log.Errorf("HTTP request error: %v", err)
		return false
	}

	// Add token
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		common.Log.Errorf("Unable to create http.Client: %v", err)
		return false
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		common.Log.Errorf("Failed to read body when pushing command response to %s.", response.Channel)
		return false
	}

	if resp.StatusCode != http.StatusOK {
		common.Log.Errorf("Failed to push response to %s. Server sent: %s.", response.Channel, body)
		return false
	}

	common.Log.Infof("Command response pushed to %s.", response.Channel)
	return true
}

func consumeCommandResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		response := rabbitmq.CommandResponse{}
		err := json.Unmarshal(d.Body, &response)
		if err != nil {
			common.Log.Errorf("Failed to decode command response : %v", err)
			d.Nack(false, false)
			return
		}

		common.Log.Infof("[cid=%s] Received message (id=%s) with type %s", d.CorrelationId, d.MessageId, d.Type)

		if !pushCommandResponse(&response) {
			d.Nack(false, false)
		} else {
			d.Ack(false)
		}
	}
}
