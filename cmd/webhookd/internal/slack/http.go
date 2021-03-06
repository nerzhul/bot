package slack

import (
	"fmt"
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	myrabbitmq "gitlab.com/nerzhul/bot/cmd/webhookd/internal/rabbitmq"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"net/http"
)

// swagger:parameters getSlackCommandResult
type slackCommandRequest struct {
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

// swagger:response slackCommandResponse
type slackCommandResponse struct {
	// in: body
	Body struct {
		// Response type
		// required: true
		ResponseType string `json:"response_type"`
		// The response content
		// required: true
		Text string `json:"text"`
	}
}

// V1ApiSlackCommand handle slack commands through HTTP
// swagger:route POST /v1/slack/commands slack-command getSlackCommandResult
//
// Handle slack commands
//
// Security:
//    jwtToken: read
//
// Responses:
//    200: slackCommandResponse
//    400: errorResponse
//    403: errorResponse
//    500: errorResponse
func V1ApiSlackCommand(c echo.Context) error {
	mcr := new(slackCommandRequest)
	if err := c.Bind(mcr); err != nil {
		common.Log.Errorf("Malformed request sent from %s, refusing slack command", c.RealIP())
		var e common.ErrorResponse
		e.Body.Message = "Bad request"
		return c.JSON(http.StatusBadRequest, e.Body)
	}

	if !common.GConfig.IsSlackTokenAllowed(mcr.Token) {
		common.Log.Errorf("Invalid token sent from %s, refusing slack command", c.RealIP())
		var e common.ErrorResponse
		e.Body.Message = "Forbidden"
		return c.JSON(http.StatusForbidden, e.Body)
	}

	if !myrabbitmq.AsyncClient.VerifyPublisher() {
		common.Log.Error("Failed to verify publisher, no command sent to broker")
		var e common.ErrorResponse
		e.Body.Message = "Server error"
		return c.JSON(http.StatusInternalServerError, e.Body)
	}

	if !myrabbitmq.AsyncClient.VerifyConsumer() {
		common.Log.Error("Failed to verify consumer, no command sent to broker")
		var e common.ErrorResponse
		e.Body.Message = "Server error"
		return c.JSON(http.StatusInternalServerError, e.Body)
	}

	consumerCfg := common.GConfig.RabbitMQ.GetConsumer("webhook")
	if consumerCfg == nil {
		common.Log.Fatalf("RabbitMQ consumer configuration 'webhook' not found, aborting.")
	}

	event := rabbitmq.CommandEvent{
		Command: fmt.Sprintf("%s %s", mcr.Command[1:], mcr.Text),
		Channel: mcr.ResponseURL,
		User:    mcr.UserName,
	}

	myrabbitmq.AsyncClient.PublishCommand(&event, consumerCfg.RoutingKey)
	return c.String(http.StatusOK, "")
}
