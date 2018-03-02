package rabbitmq

import (
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

// Publisher global rabbitmq publisher
var Publisher *rabbitmq.EventPublisher

// VerifyPublisher ensure publisher exists
func VerifyPublisher() bool {
	if Publisher == nil {
		Publisher = rabbitmq.NewEventPublisher(common.Log, &common.GConfig.RabbitMQ)
		if !Publisher.Init() {
			Publisher = nil
		}
	}

	return Publisher != nil
}
