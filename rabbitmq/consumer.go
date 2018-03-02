package rabbitmq

import (
	"github.com/op/go-logging"
	"github.com/streadway/amqp"
)

// EventConsumer publication object
type EventConsumer struct {
	conn             *amqp.Connection
	channel          *amqp.Channel
	log              *logging.Logger
	config           *RabbitMQConfig
	IncomingMessages chan amqp.Delivery
}

// NewEventConsumer creates a new EventPublisher with config & logger
func NewEventConsumer(logger *logging.Logger, config *RabbitMQConfig) *EventConsumer {
	return &EventConsumer{
		log:              logger,
		config:           config,
		IncomingMessages: make(chan amqp.Delivery),
	}
}

// Init initialize event consumer
func (ep *EventConsumer) Init() bool {
	var err error
	ep.conn, err = amqp.Dial(ep.config.URL)
	if err != nil {
		ep.log.Errorf("Failed to connect to RabbitMQ: %s", err)
		return false
	}

	ep.log.Infof("Connected to RabbitMQ on %s", ep.config.URL)

	ep.channel, err = ep.conn.Channel()
	if err != nil {
		ep.log.Errorf("Failed to open a channel: %s", err)
		return false
	}

	ep.log.Infof("RabbitMQ channel opened on %s", ep.config.URL)

	return true
}

// DeclareExchange declare exchange on event consumer
func (ep *EventConsumer) DeclareExchange(name string, durable bool) bool {
	if ep.channel == nil {
		ep.log.Fatalf("Implementation error: Consumer channel is nil")
		return false
	}

	err := ep.channel.ExchangeDeclare(
		name,
		"direct",
		durable,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		ep.log.Errorf("Failed to declare exchange %s: %s", name, err.Error())
		return false
	}

	ep.log.Infof("Exchange %s declared on RabbitMQ", name)

	return true
}

// DeclareQueue declare queue with name
func (ep *EventConsumer) DeclareQueue(name string) bool {
	if ep.channel == nil {
		ep.log.Fatalf("Implementation error: Consumer channel is nil")
		return false
	}

	_, err := ep.channel.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when usused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	if err != nil {
		ep.log.Errorf("Failed to declare queue %s: %s", name, err.Error())
		return false
	}

	ep.log.Infof("Queue %s declared on RabbitMQ", name)

	return true
}

// BindExchange bind exchange with queue
func (ep *EventConsumer) BindExchange(queue string, exchange string, routingKey string) bool {
	ep.log.Infof("Binding queue %s with exchange %s using routing key %s", queue, exchange, routingKey)
	err := ep.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		false,
		nil)

	if err != nil {
		ep.log.Errorf("Failed to bind exchange %s with queue %s: %s", exchange, queue, err.Error())
		return false
	}

	ep.log.Infof("Exchange %s bound with queue %s using routing key %s on RabbitMQ", exchange, queue, routingKey)
	return true
}

// ConsumeCallback callback function called on consuming
type ConsumeCallback func(<-chan amqp.Delivery)

// Consume consume events on queue
func (ep *EventConsumer) Consume(queue string, consumerID string, cb ConsumeCallback, autoAck bool) bool {
	if cb == nil {
		ep.log.Fatalf("ConsumeCallback is nil!")
		return false
	}

	msgs, err := ep.channel.Consume(
		queue,      // queue
		consumerID, // consumer
		autoAck,    // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		ep.log.Errorf("Failed to consume on queue %s: %s", queue, err.Error())
		return false
	}

	go cb(msgs)

	ep.log.Infof("Start consuming on queue %s", queue)

	return true
}
