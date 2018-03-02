package internal

import (
	"gitlab.com/nerzhul/bot/rabbitmq"
)

var rabbitmqPublisher *rabbitmq.EventPublisher

func verifyPublisher() bool {
	if rabbitmqPublisher == nil {
		rabbitmqPublisher = rabbitmq.NewEventPublisher(log, &gconfig.RabbitMQ)
		if !rabbitmqPublisher.Init() {
			rabbitmqPublisher = nil
		}
	}

	return rabbitmqPublisher != nil
}
