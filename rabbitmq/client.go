package rabbitmq

import (
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

// Client the RabbitMQ publisher & consumer client linked with configuration
type Client struct {
	Publisher         *EventPublisher
	Consumer          *EventConsumer
	config            *Config
	logger            *logging.Logger
	consumerNames     []string
	consumingCallback ConsumeCallback
}

// NewClient create a new RabbitMQ client
func NewClient(logger *logging.Logger, config *Config, ccb ConsumeCallback) *Client {
	return &Client{
		config:            config,
		logger:            logger,
		consumingCallback: ccb,
	}
}

// AddConsumerName add consumer name registration for consumer verifications
func (rc *Client) AddConsumerName(name string) {
	rc.consumerNames = append(rc.consumerNames, name)
}

// VerifyPublisher ensure publisher is properly created, else instantiate it
func (rc *Client) VerifyPublisher() bool {
	if rc.Publisher == nil {
		rc.Publisher = NewEventPublisher(rc.logger, rc.config)
		if !rc.Publisher.Init() {
			rc.Publisher = nil
		}
	}

	return rc.Publisher != nil
}

// PublishCommand publish a command to RabbitMQ, and hope somebody we reply on the replyTo queue
func (rc *Client) PublishCommand(cc *CommandEvent, replyTo string) bool {
	return rc.Publisher.Publish(cc, "command",
		&EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ReplyTo:       replyTo,
			ExpirationMs:  300000,
			RoutingKey:    "chat-command",
		},
	)
}

// PublishIRCCommand publish an IRC command to RabbitMQ, and hope somebody we reply on the replyTo queue
func (rc *Client) PublishIRCCommand(cc *IRCCommand, replyTo string) bool {
	return rc.Publisher.Publish(cc, "irc-command",
		&EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ReplyTo:       replyTo,
			ExpirationMs:  300000,
			RoutingKey:    "irc-command",
		},
	)
}

// PublishGitlabEvent publish incoming gitlab event to exchange
func (rc *Client) PublishGitlabEvent(event *CommandResponse) bool {
	return rc.Publisher.Publish(event, "gitlab-event",
		&EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  300000,
		},
	)
}

// VerifyConsumer ensure consumer is properly created, else instantiate it
func (rc *Client) VerifyConsumer() bool {
	if rc.consumingCallback == nil {
		rc.logger.Fatal("No consumer callback is defined, aborting.")
	}

	if rc.Consumer == nil {
		rc.Consumer = NewEventConsumer(rc.logger, rc.config)
		if !rc.Consumer.Init() {
			rc.Consumer = nil
			return false
		}

		for _, consumerName := range rc.consumerNames {
			consumerCfg := rc.config.GetConsumer(consumerName)
			if consumerCfg == nil {
				rc.logger.Fatalf("RabbitMQ consumer configuration '%s' not found, aborting.", consumerName)
			}

			if !rc.Consumer.DeclareQueue(consumerCfg.Queue) {
				rc.logger.Errorf("Failed to declare queue %s, dropping consumer", consumerCfg.Queue)
				rc.Consumer = nil
				return false
			}

			if !rc.Consumer.DeclareExchange(consumerCfg.Exchange, consumerCfg.ExchangeDurable) {
				rc.logger.Errorf("Failed to declare exchange %s, dropping consumer", consumerCfg.Exchange)
				rc.Consumer = nil
				return false
			}

			if !rc.Consumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
				rc.logger.Errorf("Failed to bind exchange %s with queue %s, dropping consumer",
					consumerCfg.Exchange, consumerCfg.Queue)
				rc.Consumer = nil
				return false
			}

			if !rc.Consumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, rc.consumingCallback, false) {
				rc.logger.Errorf("Failed to start consuming on queue %s, dropping consumer", consumerCfg.Queue)
				rc.Consumer = nil
				return false
			}
		}
	}

	return rc.Consumer != nil
}
