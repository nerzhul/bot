package internal

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

func consumeCommandQueries(msg *amqp.Delivery) {
	query := rabbitmq.CommandEvent{}
	err := json.Unmarshal(msg.Body, &query)
	if err != nil {
		log.Errorf("Failed to decode command event : %v", err)
	}

	if router == nil {
		router = &commandRouter{}
		router.init()
	}

	// Consume command queries
	if router.handleCommand(&query, msg.CorrelationId, msg.ReplyTo) {
		msg.Ack(false)
	} else {
		msg.Nack(false, true)
	}
}
