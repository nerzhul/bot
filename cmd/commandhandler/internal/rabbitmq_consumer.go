package internal

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

func consumeCommandQueries(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		query := rabbitmq.CommandEvent{}
		err := json.Unmarshal(d.Body, &query)
		if err != nil {
			log.Errorf("Failed to decode command event : %v", err)
		}

		if router == nil {
			router = &commandRouter{}
			router.init()
		}

		// Consume command queries
		if router.handleCommand(&query, d.CorrelationId, d.ReplyTo) {
			d.Ack(false)
		} else {
			d.Nack(false, true)
		}
	}
}
