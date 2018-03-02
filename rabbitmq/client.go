package rabbitmq

import "github.com/op/go-logging"

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
