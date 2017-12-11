package internal

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

type rabbitMQPublisherConfig struct {
	URL             string `yaml:"url"`
	EventExchange   string `yaml:"exchange"`
	EventRoutingKey string `yaml:"routing-key"`
}

type gitlabEventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type gitlabRabbitMQEvent struct {
	Message     string `json:"message"`
	Channel     string `json:"channel"`
	User        string `json:"user"`
	MessageType string `json:"message_type"`
}

func (gre *gitlabRabbitMQEvent) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(gre)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

func newGitlabEventPublisher() *gitlabEventPublisher {
	return &gitlabEventPublisher{}
}

func (aep *gitlabEventPublisher) init() bool {
	var err error
	aep.conn, err = amqp.Dial(gconfig.RabbitMQ.URL)
	if err != nil {
		log.Errorf("Failed to connect to RabbitMQ: %s", err)
		return false
	}

	log.Infof("Connected to RabbitMQ on %s", gconfig.RabbitMQ.URL)

	aep.channel, err = aep.conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %s", err)
		return false
	}

	log.Infof("RabbitMQ channel opened on %s", gconfig.RabbitMQ.URL)

	err = aep.channel.ExchangeDeclare(
		gconfig.RabbitMQ.EventExchange,
		"direct",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("Failed to declare exchange %s: %s", gconfig.RabbitMQ.EventExchange, err)
		return false
	}

	log.Infof("RabbitMQ exchange %s created on %s", gconfig.RabbitMQ.EventExchange, gconfig.RabbitMQ.URL)

	return true
}

func (aep *gitlabEventPublisher) publish(event *gitlabRabbitMQEvent, correlationID string) bool {
	if len(correlationID) == 0 {
		log.Fatalf("Cannot send achievement event with empty CorrelationId, aborting.")
		return false
	}

	jsonEvent, err := event.ToJSON()
	if err != nil {
		log.Errorf("Failed to publish GitlabEvent to string. Cannot publish to exchange '%s'",
			gconfig.RabbitMQ.EventExchange)
		return false
	}

	toPublish := amqp.Publishing{
		Body:          jsonEvent,
		ContentType:   "application/json",
		MessageId:     uuid.NewV4().String(),
		Timestamp:     time.Now(),
		CorrelationId: correlationID,
		Expiration:    "300000",
		Type:          "gitlab-event",
	}

	err = aep.channel.Publish(
		gconfig.RabbitMQ.EventExchange,   // exchange
		gconfig.RabbitMQ.EventRoutingKey, // routing key
		true,  // mandatory
		false, // immediate
		toPublish,
	)

	if err != nil {
		log.Errorf("Failed to publish message to exchange %s: %s", gconfig.RabbitMQ.EventExchange, err)
		// Drop rabbitmq client for a future reconnection
		rabbitmqPublisher = nil
		return false
	}

	log.Infof("[cid=%s] Message published to exchange %s", correlationID, gconfig.RabbitMQ.EventExchange)

	return true
}

func verifyPublisher() bool {
	if rabbitmqPublisher == nil {
		rabbitmqPublisher = newGitlabEventPublisher()
		if !rabbitmqPublisher.init() {
			rabbitmqPublisher = nil
		}
	}

	return rabbitmqPublisher != nil
}
