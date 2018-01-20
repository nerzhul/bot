package rabbitmq

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"io/ioutil"
	"net/http"
)

// Consumer rabbitmq consummer
var Consumer *bot.EventConsumer

func pushCommandResponse(response *bot.CommandResponse) bool {
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
		response := bot.CommandResponse{}
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

// VerifyConsumer verify and re-create rabbitmq connection if needed
func VerifyConsumer() bool {
	if Consumer == nil {
		Consumer = bot.NewEventConsumer(common.Log, &common.GConfig.RabbitMQ)
		if !Consumer.Init() {
			Consumer = nil
			return false
		}

		consumerCfg := common.GConfig.RabbitMQ.GetConsumer("webhook")
		if consumerCfg == nil {
			common.Log.Fatalf("RabbitMQ consumer configuration 'webhook' not found, aborting.")
		}

		if !Consumer.DeclareQueue(consumerCfg.Queue) {
			Consumer = nil
			return false
		}

		if !Consumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
			Consumer = nil
			return false
		}

		if !Consumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, consumeCommandResponses, false) {
			Consumer = nil
			return false
		}
	}

	return Consumer != nil
}
