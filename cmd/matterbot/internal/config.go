package internal

import (
	"gitlab.com/nerzhul/bot/rabbitmq"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	RabbitMQ rabbitmq.Config `yaml:"rabbitmq"`

	Mattermost struct {
		URL                 string   `yaml:"url"`
		WsURL               string   `yaml:"ws-url"`
		IRCWebhookURL       string   `yaml:"irc-webhook-url"`
		IRCAllowedSenders   []string `yaml:"irc-allowed-senders"`
		IRCSenderRoutingKey string   `yaml:"irc-sender-routing-key"`
		Email               string   `yaml:"email"`
		Password            string   `yaml:"password"`
		Username            string   `yaml:"username"`
		Userfirst           string   `yaml:"user-first"`
		Userlast            string   `yaml:"user-last"`
		Team                string   `yaml:"team"`
		TwitterChannel      string   `yaml:"twitter-channel"`
	} `yaml:"mattermost"`
}

var gconfig config

func (c *config) loadDefaultConfiguration() {
	c.RabbitMQ.URL = "amqp://guest:guest@localhost:5672/"
	c.RabbitMQ.EventExchange = "commands"
	c.RabbitMQ.PublisherRoutingKey = "chat-command"

	c.RabbitMQ.Consumers = map[string]rabbitmq.Consumer{
		"commands": {
			RoutingKey:      "matterbot",
			ConsumerID:      "matterbot/commands",
			Queue:           "matterbot/commands",
			Exchange:        "commands",
			ExchangeDurable: false,
		},
		"irc": {
			RoutingKey:      "irc-chat",
			ConsumerID:      "matterbot/irc",
			Queue:           "matterbot/irc",
			Exchange:        "commands",
			ExchangeDurable: false,
		},
		"twitter": {
			RoutingKey:      "matterbot",
			ConsumerID:      "matterbot/twitter",
			Queue:           "matterbot/twitter",
			Exchange:        "twitter",
			ExchangeDurable: false,
		},
	}

	c.Mattermost.URL = "http://localhost:8065"
	c.Mattermost.WsURL = "ws://localhost:8065"
	c.Mattermost.IRCWebhookURL = "http://localhost:8065/hooks/blah"
	c.Mattermost.Username = "bot"
	c.Mattermost.Password = "password"
	c.Mattermost.Email = "bot@bot.local"
	c.Mattermost.Userfirst = "Bot"
	c.Mattermost.Userlast = "Bot"
	c.Mattermost.Team = "MyTeam"
	c.Mattermost.TwitterChannel = "twitter"
	c.Mattermost.IRCSenderRoutingKey = "irc-chat-send"
}

func (c *config) isAllowedIRCSender(name string) bool {
	for _, allowedSender := range c.Mattermost.IRCAllowedSenders {
		if allowedSender == name {
			return true
		}
	}

	return false
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
