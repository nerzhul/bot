package rabbitmq

import (
	"gitlab.com/nerzhul/bot"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
)

// Publisher global rabbitmq publisher
var Publisher *bot.EventPublisher

// VerifyPublisher ensure publisher exists
func VerifyPublisher() bool {
	if Publisher == nil {
		Publisher = bot.NewEventPublisher(common.Log, &common.GConfig.RabbitMQ)
		if !Publisher.Init() {
			Publisher = nil
		}
	}

	return Publisher != nil
}
