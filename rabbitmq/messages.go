package rabbitmq

import (
	"encoding/json"
	"fmt"
)

// CommandEvent event sent to command handler
type CommandEvent struct {
	Command string `json:"command"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

// ToJSON converts CommandEvent to JSON
func (ce *CommandEvent) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(ce)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// IRCChatEvent event sent when a chat message arrives on a channel
type IRCChatEvent struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

// ToJSON converts IRCChatEvent to JSON
func (ice *IRCChatEvent) ToJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(ice)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

// IRCCommand event sent when a chat message arrives on a channel
type IRCCommand struct {
	Command string `json:"command"`
	Arg1    string `json:"arg1"`
	Arg2    string `json:"arg2"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

// ToJSON converts IRCCommand to JSON
func (ice *IRCCommand) ToJSON() ([]byte, error) {
	if ice.Command != "join" && ice.Command != "leave" && ice.Command != "list" {
		return nil, fmt.Errorf("Invalid IRCCommand command field: %s, marshaling failed.\n", ice.Command)
	}

	jsonStr, err := json.Marshal(ice)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

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
	UserScreenName string `json:"user_screen_name"`
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
