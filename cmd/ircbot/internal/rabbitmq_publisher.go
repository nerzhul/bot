package internal

import (
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

type rabbitmqClient struct {
	*rabbitmq.Client
}

func newRabbitMQClient() *rabbitmqClient {
	rc := &rabbitmqClient{}
	rc.Client = rabbitmq.NewClient(log, &gconfig.RabbitMQ)
	return rc
}

func (rc *rabbitmqClient) publishChatEvent(ce *rabbitmq.IRCChatEvent) bool {
	return rc.Publisher.Publish(ce, "irc-chat",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  1800000,
			RoutingKey:    "irc-chat",
		},
	)
}
