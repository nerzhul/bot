package internal

import (
	"fmt"
	"gitlab.com/nerzhul/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

type ircChannelConfig struct {
	Name           string `yaml:"name"`
	Password       string `yaml:"string"`
	AnswerCommands bool   `yaml:"answer-commands"`
	Hello          bool   `yaml:"hello"`
}

type config struct {
	RabbitMQ bot.RabbitMQConfig `yaml:"rabbitmq"`

	IRC struct {
		Server   string             `yaml:"server"`
		Port     uint16             `yaml:"port"`
		SSL      bool               `yaml:"ssl"`
		Name     string             `yaml:"name"`
		Password string             `yaml:"password"`
		Channels []ircChannelConfig `yaml:"channels"`
	} `yaml:"irc"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "commands"
	c.RabbitMQ.PublisherRoutingKey = "chat-command"
	c.RabbitMQ.Consumers = map[string]bot.RabbitMQConsumer{
		"ircbot": {
			RoutingKey:      "ircbot",
			ConsumerID:      "ircbot",
			Queue:           "ircbot",
			Exchange:        "commands",
			ExchangeDurable: false,
		},
	}

	c.IRC.Server = "chat.freenode.net"
	c.IRC.Port = 6697
	c.IRC.SSL = true
	c.IRC.Name = fmt.Sprintf("ircbot%d", time.Now())
}

func (c *config) getIRCChannelConfig(name string) *ircChannelConfig {
	for _, channelCfg := range c.IRC.Channels {
		if channelCfg.Name != name {
			continue
		}

		return &channelCfg
	}

	return nil
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
