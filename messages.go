package bot

import "encoding/json"

// CommandResponse command response received on RabbitMQ from command handler
type CommandResponse struct {
	Channel     string `json:"channel"`
	Message     string `json:"message"`
	User        string `json:"user"`
	MessageType string `json:"message_type"`
}

func (gre *CommandResponse) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(gre)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}
