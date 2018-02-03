package internal

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
	"strings"
)

var rabbitmqConsumer *bot.EventConsumer

func consumeResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		if d.Type == "irc-chat" {
			consumeIRCChatMessage(&d)
		} else {
			consumeCommandResponse(&d)
		}
	}
}

func consumeIRCChatMessage(msg *amqp.Delivery) {
	log.Debugf("Received message to send on IRC: %v", msg.Body)
	msg.Ack(false)
}

func consumeCommandResponse(msg *amqp.Delivery) {
	response := bot.CommandResponse{}
	err := json.Unmarshal(msg.Body, &response)
	if err != nil {
		log.Errorf("Failed to decode command response : %v", err)
	}

	if ircConn == nil {
		msg.Nack(false, true)
	}

	for _, msg := range strings.Split(response.Message, "\n") {
		if response.MessageType == "notice" {
			ircConn.Notice(response.Channel, msg)
		} else {
			ircConn.Privmsg(response.Channel, msg)
		}
	}

	msg.Ack(false)
}

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = bot.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		for _, consumerName := range []string{"ircbot", "chat"} {
			consumerCfg := gconfig.RabbitMQ.GetConsumer(consumerName)
			if consumerCfg == nil {
				log.Fatalf("RabbitMQ consumer configuration '%s' not found, aborting.", consumerName)
			}

			if !rabbitmqConsumer.DeclareQueue(consumerCfg.Queue) {
				rabbitmqConsumer = nil
				return false
			}

			if !rabbitmqConsumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
				rabbitmqConsumer = nil
				return false
			}

			if !rabbitmqConsumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, consumeResponses, false) {
				rabbitmqConsumer = nil
				return false
			}
		}
	}

	return rabbitmqConsumer != nil
}
