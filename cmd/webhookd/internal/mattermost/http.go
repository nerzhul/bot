package mattermost

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/rabbitmq"
	"net/http"
)

// swagger:parameters getMattermostCommandResult
type mattermostCommandRequest struct {
	// User ID
	//
	// in: body
	// required: true
	ChanneID    string `json:"channel_id" form:"channel_id" query:"channel_id"`
	ChannelName string `json:"channel_name" form:"channel_name" query:"channel_name"`
	Command     string `json:"command" form:"command" query:"command"`
	ResponseURL string `json:"response_url" form:"response_url" query:"response_url"`
	TeamDomain  string `json:"team_domain" form:"team_domain" query:"team_domain"`
	TeamID      string `json:"team_id" form:"team_id" query:"team_id"`
	Text        string `json:"text" form:"text" query:"text"`
	Token       string `json:"token" form:"token" query:"token"`
	UserID      string `json:"user_id" form:"user_id" query:"user_id"`
	UserName    string `json:"user_name" form:"user_name" query:"user_name"`
}

// swagger:response mattermostCommandResponse
type mattermostCommandResponse struct {
	// in: body
	Body struct {
		// Response type
		// required: true
		ResponseType string `json:"response_type"`
		// The response content
		// required: true
		Text string `json:"text"`
		// Overwrite username
		// required: false
		Username string `json:"username"`
	}
}

// V1ApiMattermostCommand handle mattermost commands through HTTP
// swagger:route POST /v1/mattermost/commands mattermost-command getMattermostCommandResult
//
// Handle mattermost commands
//
// Security:
//    jwtToken: read
//
// Responses:
//    200: mattermostCommandResponse
//    400: errorResponse
//    403: errorResponse
//    500: errorResponse
func V1ApiMattermostCommand(c echo.Context) error {
	mcr := new(mattermostCommandRequest)
	if err := c.Bind(mcr); err != nil {
		common.Log.Errorf("Malformed request sent from %s, refusing Mattermost command", c.RealIP())
		var e common.ErrorResponse
		e.Body.Message = "Bad request"
		return c.JSON(http.StatusBadRequest, e.Body)
	}

	if mcr.Token != common.GConfig.Mattermost.Token {
		common.Log.Errorf("Invalid token sent from %s, refusing Mattermost command", c.RealIP())
		var e common.ErrorResponse
		e.Body.Message = "Forbidden"
		return c.JSON(http.StatusForbidden, e.Body)
	}

	if !rabbitmq.VerifyPublisher() {
		common.Log.Error("Failed to verify publisher, no command sent to broker")
		var e common.ErrorResponse
		e.Body.Message = "Server error"
		return c.JSON(http.StatusInternalServerError, e.Body)
	}

	if !rabbitmq.VerifyPublisher() {
		common.Log.Error("Failed to verify consumer, no command sent to broker")
		var e common.ErrorResponse
		e.Body.Message = "Server error"
		return c.JSON(http.StatusInternalServerError, e.Body)
	}

	consumerCfg := common.GConfig.RabbitMQ.GetConsumer("webhook")
	if consumerCfg == nil {
		common.Log.Fatalf("RabbitMQ consumer configuration 'webhook' not found, aborting.")
	}

	event := bot.CommandEvent{
		Command: fmt.Sprintf("%s %s", mcr.Command[1:], mcr.Text),
		Channel: mcr.ResponseURL,
		User:    mcr.UserName,
	}

	rabbitmq.Publisher.Publish(
		&event,
		"command",
		&bot.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ReplyTo:       consumerCfg.RoutingKey,
			ExpirationMs:  300000,
			RoutingKey:    "chat-command",
		},
	)
	return c.JSON(http.StatusOK, nil)
}
