package internal

import (
	"gitlab.com/nerzhul/bot/rabbitmq"
)

type rabbitmqClient struct {
	*rabbitmq.Client
}

func newRabbitMQClient() *rabbitmqClient {
	rc := &rabbitmqClient{}
	rc.Client = rabbitmq.NewClient(log, &gconfig.RabbitMQ, consumeCommandQueries)
	return rc
}
