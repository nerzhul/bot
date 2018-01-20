package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
)

// Consumer rabbitmq consummer
var Consumer *bot.EventConsumer

func consumeCommandResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		response := bot.CommandResponse{}
		err := json.Unmarshal(d.Body, &response)
		if err != nil {
			common.Log.Errorf("Failed to decode command response : %v", err)
		}

		common.Log.Debugf("command response %v", response)
		if response.MessageType == "notice" {
			// ircConn.Notice(response.Channel, msg)
		} else {
			// ircConn.Privmsg(response.Channel, msg)
		}

		d.Ack(false)
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
