package internal

import (
	"gitlab.com/nerzhul/gitlab-hook"
)

var rabbitmqConsumer *bot.EventConsumer

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

		if !rabbitmqConsumer.Consume(gconfig.RabbitMQ.EventQueue) {
			rabbitmqConsumer = nil
			return false
		}
	}

	return rabbitmqConsumer != nil
}
