package rabbitmq

import (
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

var AsyncClient *rabbitmqClient

type rabbitmqClient struct {
	*rabbitmq.Client
}

func NewRabbitMQClient() *rabbitmqClient {
	rc := &rabbitmqClient{}
	rc.Client = rabbitmq.NewClient(common.Log, &common.GConfig.RabbitMQ, consumeCommandResponses)
	return rc
}

func (rc *rabbitmqClient) PublishGitlabEvent(event *rabbitmq.CommandResponse) {
	rc.Publisher.Publish(event, "gitlab-event",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  300000,
		},
	)
}
