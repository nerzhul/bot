package rabbitmq

import (
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

// AsyncClient the async client
var AsyncClient *MyAsyncClient

// MyAsyncClient the private webhook rabbit client
type MyAsyncClient struct {
	*rabbitmq.Client
}

// NewRabbitMQClient create asynchronous client
func NewRabbitMQClient() *MyAsyncClient {
	rc := &MyAsyncClient{}
	rc.Client = rabbitmq.NewClient(common.Log, &common.GConfig.RabbitMQ, consumeCommandResponses)
	return rc
}
