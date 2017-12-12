package internal

import (
	"encoding/json"
	"gitlab.com/nerzhul/gitlab-hook"
)

var rabbitmqPublisher *bot.EventPublisher

type commandEvent struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

func (ce *commandEvent) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(ce)
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
