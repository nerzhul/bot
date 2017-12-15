package internal

import (
	"gitlab.com/nerzhul/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ bot.RabbitMQConfig `yaml:"rabbitmq"`

	Twitter struct {
		ConsumerKey    string `yaml:"consumer-key"`
		ConsumerSecret string `yaml:"consumer-secret"`
		Token          string `yaml:"token"`
		TokenSecret    string `yaml:"token-secret"`
	} `yaml:"twitter"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "commands"
	c.RabbitMQ.PublisherRoutingKey = "slackbot"
	c.RabbitMQ.Consumers = map[string]bot.RabbitMQConsumer{
		"twitterbot": {
			RoutingKey:      "twitterbot",
			ConsumerID:      "twitterbot",
			Queue:           "twitterbot",
			Exchange:        "commands",
			ExchangeDurable: false,
		},
	}

	c.Twitter.ConsumerKey = "consumer-key"
	c.Twitter.ConsumerSecret = "consumer-secret"
	c.Twitter.Token = "token"
	c.Twitter.TokenSecret = "token-secret"
}

func loadConfiguration(path string) {
	gconfig.loadDefaultConfiguration()

	if len(path) == 0 {
		log.Info("Configuration path is empty, using default configuration.")
		return
	}

	log.Infof("Loading configuration from '%s'...", path)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read YAML file #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &gconfig)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	log.Infof("Configuration loaded from '%s'.", path)
}
