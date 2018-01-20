package slack

import (
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
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

	if mcr.Token != common.GConfig.Slack.Token {
		common.Log.Errorf("Invalid token sent from %s, refusing slack command", c.RealIP())
		var e common.ErrorResponse
		e.Body.Message = "Forbidden"
		return c.JSON(http.StatusForbidden, e.Body)
	}

	mcrp := slackCommandResponse{}
	mcrp.Body.ResponseType = "ephemeral"
	mcrp.Body.Text = "This is a test"
	return c.JSON(http.StatusOK, mcrp.Body)
}