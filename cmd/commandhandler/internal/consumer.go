package internal

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
)

var rabbitmqConsumer *bot.EventConsumer

func consumeCommandQueries(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		query := bot.CommandEvent{}
		err := json.Unmarshal(d.Body, &query)
		if err != nil {
			log.Errorf("Failed to decode command event : %v", err)
		}

		// Consume command queries

		d.Ack(false)
	}
}

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = bot.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		consumerCfg := gconfig.RabbitMQ.GetConsumer("commandhandler")
		if consumerCfg == nil {
			log.Fatalf("RabbitMQ consumer configuration 'commandhandler' not found, aborting.")
		}

		if !rabbitmqConsumer.DeclareQueue(consumerCfg.Queue) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, consumeCommandQueries, false) {
			rabbitmqConsumer = nil
			return false
		}
	}

	return rabbitmqConsumer != nil
}
