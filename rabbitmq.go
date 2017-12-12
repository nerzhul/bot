package bot

// RabbitMQConfig standard publisher configuration
type RabbitMQConfig struct {
	URL                 string `yaml:"url"`
	EventExchange       string `yaml:"exchange"`
	PublisherRoutingKey string `yaml:"publisher-routing-key"`
	ConsumerID          string `yaml:"consumer-id"`
	ConsumerRoutingKey  string `yaml:"consumer-routing-key"`
	EventQueue          string `yaml:"queue"`
}
