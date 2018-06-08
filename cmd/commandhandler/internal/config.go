package internal

import (
	"gitlab.com/nerzhul/bot/rabbitmq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`
	Scaleway struct {
		URL           string `yaml:"url"`
		Token         string `yaml:"token"`
		BuildServerID string `yaml:"build-server-id"`
	} `yaml:"scaleway"`
	OpenWeatherMap struct {
		APIKey string `yaml:"apikey"`
		Lang   string `yaml:"lang"`
		Unit   string `yaml:"unit"`
	} `yaml:"openweathermap"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "commands"
	c.RabbitMQ.EventExchangeType = "direct"
	c.RabbitMQ.PublisherRoutingKey = ""
	c.RabbitMQ.Consumers = map[string]rabbitmq.Consumer{
		"commandhandler": {
			RoutingKey:      "chat-command",
			ConsumerID:      "botcommand",
			Queue:           "botcommand.direct",
			Exchange:        "commands",
			ExchangeDurable: false,
			ExchangeType:    "direct",
		},
	}
	c.Scaleway.URL = "https://cp-par1.scaleway.com"
	c.OpenWeatherMap.APIKey = ""
	c.OpenWeatherMap.Lang = "en"
	c.OpenWeatherMap.Unit = "C"
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
