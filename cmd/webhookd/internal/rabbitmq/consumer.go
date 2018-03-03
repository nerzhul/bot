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
)

func pushCommandResponse(response *rabbitmq.CommandResponse) bool {
	common.Log.Debugf("Received command response for user %s (callback %s)", response.User, response.Channel)
	pushResponse := fmt.Sprintf(`{"response_type": "ephemeral", "text": "%s"}`, response.Message)
	req, err := http.NewRequest("POST", response.Channel, bytes.NewBuffer([]byte(pushResponse)))
	if err != nil {
		common.Log.Errorf("HTTP request error: %v", err)
		return false
	}

	// Add token
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
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
		}

		if !pushCommandResponse(&response) {
			d.Nack(false, false)
		} else {
			d.Ack(false)
		}
	}
}
