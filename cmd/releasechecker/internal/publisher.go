package internal

import (
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

var rabbitmqPublisher *rabbitmq.EventPublisher

func verifyPublisher() bool {
	if rabbitmqPublisher == nil || !rabbitmqPublisher.IsValid() {
		rabbitmqPublisher = rabbitmq.NewEventPublisher(log, &gconfig.RabbitMQ)
		if !rabbitmqPublisher.Init() {
			rabbitmqPublisher = nil
		}
	}

	return rabbitmqPublisher != nil
}

func publishAnnouncement(ce *rabbitmq.AnnouncementMessage) bool {
	return rabbitmqPublisher.Publish(ce, "announcement",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  18000000,
			RoutingKey:    "announcement",
		},
	)
}
