package bot

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
	config  *RabbitMQConfig
}

// Event interface
type Event interface {
	ToJSON() ([]byte, error)
}

// NewEventPublisher creates a new EventPublisher with config & logger
func NewEventPublisher(logger *logging.Logger, config *RabbitMQConfig) *EventPublisher {
	return &EventPublisher{
		log:    logger,
		config: config,
	}
}

// Init initialize event publisher
func (ep *EventPublisher) Init() bool {
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

	err = ep.channel.ExchangeDeclare(
		ep.config.EventExchange,
		"direct",
		ep.config.EventExchangeDurable,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		ep.log.Errorf("Failed to declare exchange %s: %s", ep.config.EventExchange, err)
		return false
	}

	ep.log.Infof("RabbitMQ exchange %s created on %s", ep.config.EventExchange, ep.config.URL)

	return true
}

// Publish publish event
func (ep *EventPublisher) Publish(event Event, eventType string, correlationID string, replyTo string, expirationMs uint32) bool {
	if len(correlationID) == 0 {
		ep.log.Fatalf("Cannot send achievement event with empty CorrelationId, aborting.")
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
		CorrelationId: correlationID,
		Type:          eventType,
	}

	if expirationMs != 0 {
		toPublish.Expiration = fmt.Sprintf("%d", expirationMs)
	}

	if len(replyTo) > 0 {
		toPublish.ReplyTo = replyTo
	}

	err = ep.channel.Publish(
		ep.config.EventExchange,       // exchange
		ep.config.PublisherRoutingKey, // routing key
		true,  // mandatory
		false, // immediate
		toPublish,
	)

	if err != nil {
		ep.log.Errorf("Failed to publish message to exchange %s: %s", ep.config.EventExchange, err)
		return false
	}

	ep.log.Infof("[cid=%s] Message published to exchange %s (routing-key %s)", correlationID, ep.config.EventExchange,
		ep.config.PublisherRoutingKey)

	return true
}
