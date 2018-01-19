package mattermost

import (
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
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
		Text         string `json:"text"`
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
		return c.JSON(http.StatusBadRequest, nil)
	}

	if mcr.Token != common.GConfig.Mattermost.Token {
		var e common.ErrorResponse
		e.Body.Message = "Forbidden"
		return c.JSON(http.StatusForbidden, nil)
	}

	mcrp := mattermostCommandResponse{}
	mcrp.Body.ResponseType = "in_channel"
	mcrp.Body.Text = "This is a test"
	return c.JSON(http.StatusBadRequest, mcrp.Body)
}
