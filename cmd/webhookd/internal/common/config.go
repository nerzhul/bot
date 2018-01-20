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
		EnableHook       bool     `yaml:"enable-hook"`
		Tokens           []string `yaml:"tokens"`
		ResponseUsername string   `yaml:"response-username"`
	} `yaml:"mattermost"`

	Slack struct {
		EnableHook bool     `yaml:"enable-hook"`
		Tokens     []string `yaml:"tokens"`
	} `yaml:"slack"`
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
	c.Mattermost.Tokens = []string{}
	c.Mattermost.ResponseUsername = "webhook"

	c.Slack.EnableHook = true
	c.Slack.Tokens = []string{}
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

// IsMattermostTokenAllowed verify if mattermost token is allowed to connect
func (c *config) IsMattermostTokenAllowed(token string) bool {
	for _, t := range c.Mattermost.Tokens {
		if t == token {
			return true
		}
	}
	return false
}

// IsSlackTokenAllowed verify if mattermost token is allowed to connect
func (c *config) IsSlackTokenAllowed(token string) bool {
	for _, t := range c.Slack.Tokens {
		if t == token {
			return true
		}
	}
	return false
}
