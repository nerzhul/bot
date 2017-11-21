package internal

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ struct {
		Url        string `yaml:"url"`
		Exchange   string `yaml:"exchange"`
		RoutingKey string `yaml:"routing-key"`
	}

	Http struct {
		Port uint16 `yaml:"port"`
	}
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.Url = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.Exchange = "gitlab"
	c.RabbitMQ.RoutingKey = "gitlab/events"

	c.Http.Port = 8080
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
