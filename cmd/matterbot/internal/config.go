package internal

import (
	"gitlab.com/nerzhul/bot"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ bot.RabbitMQConfig `yaml:"rabbitmq"`

	Mattermost struct {
		URL            string `yaml:"url"`
		WsURL          string `yaml:"ws-url"`
		Email          string `yaml:"email"`
		Password       string `yaml:"password"`
		Username       string `yaml:"username"`
		Userfirst      string `yaml:"user-first"`
		Userlast       string `yaml:"user-last"`
		Team           string `yaml:"team"`
		TwitterChannel string `yaml:"twitter-channel"`
	} `yaml:"mattermost"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "commands"
	c.RabbitMQ.PublisherRoutingKey = "chat-command"

	c.RabbitMQ.Consumers = map[string]bot.RabbitMQConsumer{
		"commands": {
			RoutingKey:      "slackbot",
			ConsumerID:      "slackbot/commands",
			Queue:           "slackbot/commands",
			Exchange:        "commands",
			ExchangeDurable: false,
		},
		"twitter": {
			RoutingKey:      "slackbot",
			ConsumerID:      "slackbot/twitter",
			Queue:           "slackbot/twitter",
			Exchange:        "twitter",
			ExchangeDurable: false,
		},
	}

	c.Mattermost.URL = "http://localhost:8065"
	c.Mattermost.WsURL = "ws://localhost:8065"
	c.Mattermost.Username = "bot"
	c.Mattermost.Password = "password"
	c.Mattermost.Email = "bot@bot.local"
	c.Mattermost.Userfirst = "Bot"
	c.Mattermost.Userlast = "Bot"
	c.Mattermost.Team = "MyTeam"
	c.Mattermost.TwitterChannel = "twitter"
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