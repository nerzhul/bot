package rabbitmq

import (
	"github.com/op/go-logging"
	"github.com/satori/go.uuid"
)

// Client the RabbitMQ publisher & consumer client linked with configuration
type Client struct {
	Publisher *EventPublisher
	Consumer  *EventConsumer
	config    *Config
	logger    *logging.Logger
}

// NewClient create a new RabbitMQ client
func NewClient(logger *logging.Logger, config *Config) *Client {
	return &Client{
		config: config,
		logger: logger,
	}
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
		},
	)
}
