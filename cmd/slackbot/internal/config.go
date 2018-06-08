package internal

import (
	"gitlab.com/nerzhul/bot/rabbitmq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`

	Slack struct {
		APIKey         string `yaml:"api-key"`
		TwitterChannel string `yaml:"twitter-channel"`
	} `yaml:"slack"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "commands"
	c.RabbitMQ.EventExchangeType = "direct"
	c.RabbitMQ.PublisherRoutingKey = "chat-command"

	c.RabbitMQ.Consumers = map[string]rabbitmq.Consumer{
		"commands": {
			RoutingKey:      "slackbot",
			ConsumerID:      "slackbot/commands",
			Queue:           "slackbot/commands",
			Exchange:        "commands",
			ExchangeDurable: false,
			ExchangeType:    "direct",
		},
		"twitter": {
			RoutingKey:      "slackbot",
			ConsumerID:      "slackbot/twitter",
			Queue:           "slackbot/twitter",
			Exchange:        "twitter",
			ExchangeDurable: false,
			ExchangeType:    "direct",
		},
	}

	c.Slack.TwitterChannel = "channel-id"
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
