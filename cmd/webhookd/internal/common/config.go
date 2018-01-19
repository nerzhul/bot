package common

import (
	"gitlab.com/nerzhul/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ bot.RabbitMQConfig `yaml:"rabbitmq"`

	HTTP struct {
		Port uint16 `yaml:"port"`
	}

	GitlabProjectsMapping map[string][]string `yaml:"gitlab-projects-mapping"`

	Mattermost struct {
		EnableHook       bool   `yaml:"enable-hook"`
		Token            string `yaml:"token"`
		ResponseUsername string `yaml:"response-username"`
	} `yaml:"mattermost"`
}

// GConfig global configuration
var GConfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "gitlab"
	c.RabbitMQ.EventExchangeDurable = true
	c.RabbitMQ.PublisherRoutingKey = "gitlab/events"

	c.HTTP.Port = 8080
	c.GitlabProjectsMapping = make(map[string][]string)

	c.Mattermost.EnableHook = true
	c.Mattermost.Token = ""
	c.Mattermost.ResponseUsername = "webhook"
}

// LoadConfiguration load configuration from path
func LoadConfiguration(path string) {
	GConfig.loadDefaultConfiguration()

	if len(path) == 0 {
		Log.Info("Configuration path is empty, using default configuration.")
		return
	}

	Log.Infof("Loading configuration from '%s'...", path)

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		Log.Fatalf("Failed to read YAML file #%v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &GConfig)
	if err != nil {
		Log.Fatalf("error: %v", err)
	}

	Log.Infof("Configuration loaded from '%s'.", path)
}
