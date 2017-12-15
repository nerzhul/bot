package bot

import "encoding/json"

// CommandResponse command response received on RabbitMQ from command handler
type CommandResponse struct {
	Channel     string `json:"channel"`
	Message     string `json:"message"`
	User        string `json:"user"`
	MessageType string `json:"message_type"`
}

// ToJSON converts CommandResponse to json string
func (gre *CommandResponse) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(gre)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// TweetMessage twitter reduced message for transport on rabbitmq
type TweetMessage struct {
	Message        string `json:"message"`
	Username       string `json:"username"`
	UserScreenName string `json:"userèscreen_name"`
	Date           string `json:"date"`
}

// ToJSON converts to json
func (ce *TweetMessage) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(ce)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}
