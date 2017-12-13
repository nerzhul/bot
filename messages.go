package bot

// CommandResponse command response received on RabbitMQ from command handler
type CommandResponse struct {
	Channel string `json:"channel"`
	Message string `json:"message"`
	User    string `json:"user"`
}
