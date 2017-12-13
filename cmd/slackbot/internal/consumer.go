package internal

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/gitlab-hook"
)

var rabbitmqConsumer *bot.EventConsumer
var slackMsgID = 0

func consumeCommandResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		response := bot.CommandResponse{}
		err := json.Unmarshal(d.Body, &response)
		if err != nil {
			log.Errorf("Failed to decode command response : %v", err)
		}

		// Send message on slack
		slackMsgID++
		slackRTM.SendMessage(&slack.OutgoingMessage{
			ID:      slackMsgID,
			Type:    "message",
			Channel: response.Channel,
			Text:    response.Message,
		})
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

		if !rabbitmqConsumer.Consume(gconfig.RabbitMQ.EventQueue, consumeCommandResponses) {
			rabbitmqConsumer = nil
			return false
		}
	}

	return rabbitmqConsumer != nil
}
