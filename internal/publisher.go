package internal

import (
	"encoding/json"
	"github.com/satori/go.uuid"
	"github.com/streadway/amqp"
	"time"
)

type rabbitMQPublisherConfig struct {
	Url             string `yaml:"url"`
	EventExchange   string `yaml:"achievement-exchange"`
	EventRoutingKey string `yaml:"achievement-routing-key"`
}

type gitlabEventPublisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type gitlabRabbitMQEvent struct {
	Message  string   `json:"message"`
	Channels []string `json:"channels"`
}

func (gre *gitlabRabbitMQEvent) ToJson() ([]byte, error) {
	jsonStr, err := json.Marshal(gre)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

func NewGitlabEventPublisher() *gitlabEventPublisher {
	return &gitlabEventPublisher{}
}

func (aep *gitlabEventPublisher) Init() bool {
	var err error
	aep.conn, err = amqp.Dial(gconfig.RabbitMQ.Url)
	if err != nil {
		log.Errorf("Failed to connect to RabbitMQ: %s", err)
		return false
	}

	log.Infof("Connected to RabbitMQ on %s", gconfig.RabbitMQ.Url)

	aep.channel, err = aep.conn.Channel()
	if err != nil {
		log.Errorf("Failed to open a channel: %s", err)
		return false
	}

	log.Infof("RabbitMQ channel opened on %s", gconfig.RabbitMQ.Url)

	err = aep.channel.ExchangeDeclare(
		gconfig.RabbitMQ.EventExchange,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Errorf("Failed to declare exchange %s: %s", gconfig.RabbitMQ.EventExchange, err)
		return false
	}

	log.Infof("RabbitMQ exchange %s created on %s", gconfig.RabbitMQ.EventExchange, gconfig.RabbitMQ.Url)

	return true
}

func (aep *gitlabEventPublisher) Publish(event *gitlabRabbitMQEvent, correlationId string) bool {
	if len(correlationId) == 0 {
		log.Fatalf("Cannot send achievement event with empty CorrelationId, aborting.")
		return false
	}

	jsonEvent, err := event.ToJson()
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
		CorrelationId: correlationId,
	}

	err = aep.channel.Publish(
		gconfig.RabbitMQ.EventExchange,   // exchange
		gconfig.RabbitMQ.EventRoutingKey, // routing key
		false, // mandatory
		false, // immediate
		toPublish,
	)

	if err != nil {
		log.Errorf("Failed to publish message to exchange %s: %s", gconfig.RabbitMQ.EventExchange, err)
		// Drop rabbitmq client for a future reconnection
		rabbitmqPublisher = nil
		return false
	}

	log.Infof("[cid=%s] Message published to exchange %s", correlationId, gconfig.RabbitMQ.EventExchange)

	return true
}

func verifyPublisher() bool {
	if rabbitmqPublisher == nil {
		rabbitmqPublisher = NewGitlabEventPublisher()
		if !rabbitmqPublisher.Init() {
			rabbitmqPublisher = nil
		}
	}

	return rabbitmqPublisher != nil
}
