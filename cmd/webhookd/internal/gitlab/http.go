package gitlab

import (
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"net/http"
)

// swagger:response gitlabEventResponse
type gitlabEventResponse struct {
	// in: body
	Body struct {
		// required: true
		Status string `json:"status,required"`
	}
}

// swagger:parameters pushGitlabEvent
type gitlabEventType struct {
	// X-Gitlab-Event header
	//
	// in: header
	// required: true
	Type string
}

// V1ApiGitlabEvent Gitlab hook event
// swagger:route POST /v1/gitlab/event gitlab pushGitlabEvent
//
// Fetch user's achievements
//
// Security:
//    jwtToken: read
//
// Responses:
//    200: gitlabEventResponse
//    400: errorResponse
//    403: errorResponse
//    500: errorResponse
func V1ApiGitlabEvent(c echo.Context) error {
	eventType := gitlabEventType{Type: c.Request().Header.Get("X-Gitlab-Event")}
	switch eventType.Type {
	case "Push Hook":
		if !handleGitlabPush(c) {
			var e common.ErrorResponse
			e.Body.Message = "Internal error"
			c.JSON(http.StatusInternalServerError, e.Body)
		}
		return nil
	case "Tag Push Hook":
		if !handleGitlabTagPush(c) {
			var e common.ErrorResponse
			e.Body.Message = "Internal error"
			c.JSON(http.StatusInternalServerError, e.Body)
		}
		return nil
	case "Merge Request Hook":
		if !handleGitlabMergeRequest(c) {
			var e common.ErrorResponse
			e.Body.Message = "Internal error"
			c.JSON(http.StatusInternalServerError, e.Body)
		}
		return nil
	case "Issue Hook":
	case "Note Hook":
	case "Wiki Page Hook":
	case "Pipeline Hook":
	case "Build Hook":
		common.Log.Warningf("Unhandled X-Gitlab-Event %s", eventType.Type)
		return nil
		break
	default:
		var e common.ErrorResponse
		e.Body.Message = "Invalid request"
		c.JSON(http.StatusBadRequest, e.Body)
		return nil
	}

	return nil
}
