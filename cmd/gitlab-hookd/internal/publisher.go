package internal

import (
	"encoding/json"
	"gitlab.com/nerzhul/gitlab-hook"
)

type gitlabRabbitMQEvent struct {
	Message     string `json:"message"`
	Channel     string `json:"channel"`
	User        string `json:"user"`
	MessageType string `json:"message_type"`
}

func (gre *gitlabRabbitMQEvent) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(gre)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

func verifyPublisher() bool {
	if rabbitmqPublisher == nil {
		rabbitmqPublisher = bot.NewEventPublisher(log, &gconfig.RabbitMQ)
		if !rabbitmqPublisher.Init() {
			rabbitmqPublisher = nil
		}
	}

	return rabbitmqPublisher != nil
}
