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

		consumerCfg := gconfig.RabbitMQ.GetConsumer("ircbot")
		if consumerCfg == nil {
			log.Fatalf("RabbitMQ consumer configuration 'ircbot' not found, aborting.")
		}

		if !rabbitmqConsumer.DeclareQueue(consumerCfg.Queue) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, consumeCommandResponses, false) {
			rabbitmqConsumer = nil
			return false
		}
	}

	return rabbitmqConsumer != nil
}
