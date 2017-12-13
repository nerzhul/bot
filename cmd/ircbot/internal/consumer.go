package internal

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
	"strings"
)

var rabbitmqConsumer *bot.EventConsumer

func consumeCommandResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		response := bot.CommandResponse{}
		err := json.Unmarshal(d.Body, &response)
		if err != nil {
			log.Errorf("Failed to decode command response : %v", err)
		}

		if ircConn == nil {
			d.Nack(false, true)
		}

		for _, msg := range strings.Split(response.Message, "\n") {
			if response.MessageType == "notice" {
				ircConn.Notice(response.Channel, msg)
			} else {
				ircConn.Privmsg(response.Channel, msg)
			}
		}

		d.Ack(false)
	}
}

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = bot.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.DeclareQueue(gconfig.RabbitMQ.EventQueue) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.BindExchange(gconfig.RabbitMQ.EventQueue,
			gconfig.RabbitMQ.EventExchange, gconfig.RabbitMQ.ConsumerRoutingKey) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.Consume(gconfig.RabbitMQ.EventQueue, consumeCommandResponses, false) {
			rabbitmqConsumer = nil
			return false
		}
	}

	return rabbitmqConsumer != nil
}
