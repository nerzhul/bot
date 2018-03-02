package internal

import (
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

type rabbitmqClient struct {
	publisher *rabbitmq.EventPublisher
	consumer  *rabbitmq.EventConsumer
}

func (rc *rabbitmqClient) verifyPublisher() bool {
	if rc.publisher == nil {
		rc.publisher = rabbitmq.NewEventPublisher(log, &gconfig.RabbitMQ)
		if !rc.publisher.Init() {
			rc.publisher = nil
		}
	}

	return rc.publisher != nil
}

func (rc *rabbitmqClient) publishChatEvent(ce *rabbitmq.IRCChatEvent) bool {
	return rc.publisher.Publish(ce, "irc-chat",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  1800000,
			RoutingKey:    "irc-chat",
		},
	)
}

func (rc *rabbitmqClient) publishCommand(cc *rabbitmq.CommandEvent, replyTo string) bool {
	return rc.publisher.Publish(cc, "command",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ReplyTo:       replyTo,
			ExpirationMs:  300000,
		},
	)
}
