package bot

// RabbitMQConsumer single consumer object
type RabbitMQConsumer struct {
	ConsumerID      string `yaml:"consumer-id"`
	Queue           string `yaml:"queue"`
	Exchange        string `yaml:"exchange"`
	ExchangeDurable bool   `yaml:"exchange-durable"`
	RoutingKey      string `yaml:"routing-key"`
}

// RabbitMQConfig standard configuration
type RabbitMQConfig struct {
	URL                  string `yaml:"url"`
	EventExchange        string `yaml:"exchange"`
	EventExchangeDurable bool   `yaml:"exchange-durable"`
	PublisherRoutingKey  string `yaml:"publisher-routing-key"`
	Consumers            map[string]RabbitMQConsumer
}

// GetConsumer retrieve consumer frm RabbitMQConfig consumer list.
// nil if not found
func (c *RabbitMQConfig) GetConsumer(name string) *RabbitMQConsumer {
	if rmc, ok := c.Consumers[name]; ok {
		return &rmc
	}
	return nil
}
