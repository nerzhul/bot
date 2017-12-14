package internal

import (
	"gitlab.com/nerzhul/bot"
)

var rabbitmqPublisher *bot.EventPublisher

func verifyPublisher() bool {
	if rabbitmqPublisher == nil {
		rabbitmqPublisher = bot.NewEventPublisher(log, &gconfig.RabbitMQ)
		if !rabbitmqPublisher.Init() {
			rabbitmqPublisher = nil
		}
	}

	return rabbitmqPublisher != nil
}
