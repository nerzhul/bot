package internal

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ rabbitMQPublisherConfig `yaml:"rabbitmq"`

	HTTP struct {
		Port uint16 `yaml:"port"`
	}

	ProjectsMapping map[string][]string `yaml:"projects-mapping"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "gitlab"
	c.RabbitMQ.EventRoutingKey = "gitlab/events"

	c.HTTP.Port = 8080
	c.ProjectsMapping = make(map[string][]string)
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
