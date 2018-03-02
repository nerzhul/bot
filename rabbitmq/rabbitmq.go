package rabbitmq

// Consumer single consumer object
type Consumer struct {
	ConsumerID      string `yaml:"consumer-id"`
	Queue           string `yaml:"queue"`
	Exchange        string `yaml:"exchange"`
	ExchangeDurable bool   `yaml:"exchange-durable"`
	RoutingKey      string `yaml:"routing-key"`
}

// Config standard configuration
type Config struct {
	URL                  string `yaml:"url"`
	EventExchange        string `yaml:"exchange"`
	EventExchangeDurable bool   `yaml:"exchange-durable"`
	PublisherRoutingKey  string `yaml:"publisher-routing-key"`
	Consumers            map[string]Consumer
}

// GetConsumer retrieve consumer frm Config consumer list.
// nil if not found
func (c *Config) GetConsumer(name string) *Consumer {
	if rmc, ok := c.Consumers[name]; ok {
		return &rmc
	}
	return nil
}
