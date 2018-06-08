package rabbitmq

import (
	"fmt"
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

// EventPublisher publication object
type EventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	log     *logging.Logger
	config  *Config
	valid   bool
}

// Event interface
type Event interface {
	ToJSON() ([]byte, error)
}

// EventOptions event options on publication
type EventOptions struct {
	CorrelationID string
	ReplyTo       string
	ExpirationMs  uint32
	RoutingKey    string
}

// NewEventPublisher creates a new EventPublisher with config & logger
func NewEventPublisher(logger *logging.Logger, config *Config) *EventPublisher {
	return &EventPublisher{
		log:    logger,
		config: config,
		valid:  false,
	}
}

// Init initialize event publisher
func (ep *EventPublisher) Init() bool {
	ep.valid = false

	var err error
	ep.conn, err = amqp.Dial(ep.config.URL)
	if err != nil {
		ep.log.Errorf("Failed to connect to RabbitMQ: %s", err)
		return ep.valid
	}

	ep.log.Infof("Connected to RabbitMQ on %s", ep.config.URL)

	ep.channel, err = ep.conn.Channel()
	if err != nil {
		ep.log.Errorf("Failed to open a channel: %s", err)
		return ep.valid
	}

	ep.log.Infof("RabbitMQ channel opened on %s", ep.config.URL)

	err = ep.channel.ExchangeDeclare(
		ep.config.EventExchange,
		ep.config.EventExchangeType,
		ep.config.EventExchangeDurable,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		ep.log.Errorf("Failed to declare exchange %s: %s", ep.config.EventExchange, err)
		return ep.valid
	}

	ep.log.Infof("RabbitMQ exchange %s created on %s", ep.config.EventExchange, ep.config.URL)

	ep.valid = true
	return ep.valid
}

// IsValid return the valid flag
func (ep *EventPublisher) IsValid() bool {
	return ep.valid
}

// Publish publish event
func (ep *EventPublisher) Publish(event Event, eventType string, options *EventOptions) bool {
	if len(options.CorrelationID) == 0 {
		ep.log.Fatalf("Cannot send event with empty CorrelationId, aborting.")
		return false
	}

	jsonEvent, err := event.ToJSON()
	if err != nil {
		ep.log.Errorf("Failed to convert event to string. Cannot publish to exchange '%s'",
			ep.config.EventExchange)
		return false
	}

	toPublish := amqp.Publishing{
		Body:          jsonEvent,
		ContentType:   "application/json",
		MessageId:     uuid.NewV4().String(),
		Timestamp:     time.Now(),
		CorrelationId: options.CorrelationID,
		Type:          eventType,
	}

	if options.ExpirationMs != 0 {
		toPublish.Expiration = fmt.Sprintf("%d", options.ExpirationMs)
	}

	if len(options.ReplyTo) > 0 {
		toPublish.ReplyTo = options.ReplyTo
	}

	routingKey := options.RoutingKey
	// If routing key not set, use global options
	if len(routingKey) == 0 {
		routingKey = ep.config.PublisherRoutingKey
	}

	err = ep.channel.Publish(
		ep.config.EventExchange, // exchange
		routingKey,              // routing key
		true,                    // mandatory
		false,                   // immediate
		toPublish,
	)

	if err != nil {
		ep.log.Errorf("Failed to publish message to exchange %s: %s", ep.config.EventExchange, err)
		// When publication failed, invalidate the publisher
		ep.valid = false
		return false
	}

	ep.log.Infof("[cid=%s] Message published to exchange %s (routing-key %s)",
		options.CorrelationID, ep.config.EventExchange, routingKey)

	return true
}
