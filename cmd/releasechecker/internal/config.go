package internal

import (
	"gitlab.com/nerzhul/bot/db"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	DB          db.Config       `yaml:"database"`
	RabbitMQ    rabbitmq.Config `yaml:"rabbitmq"`
	GithubToken string          `yaml:"github-token"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() bool {
	c.DB.URL = "host=postgres dbname=releasechecker user=releasechecker password=releasechecker"
	c.DB.MaxOpenConns = 10
	c.DB.MaxIdleConns = 5
	c.DB.MigrationSource = "file:///migrations"

	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "announcements"
	c.RabbitMQ.EventExchangeType = "fanout"
	c.RabbitMQ.PublisherRoutingKey = "release"

	c.GithubToken = ""

	return true
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
