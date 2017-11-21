package internal

import (
	"github.com/labstack/echo"
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

// swagger:response errorResponse
type errorResponse struct {
	// in: body
	Body struct {
		// required: true
		Message string `json:"message,required"`
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
func v1ApiGitlabEvent(c echo.Context) error {
	eventType := gitlabEventType{Type: c.Request().Header.Get("X-Gitlab-Event")}
	switch eventType.Type {
	case "Push Hook":
	case "Tag Push Hook":
	case "Issue Hook":
	case "Note Hook":
	case "Merge Request Hook":
	case "Wiki Page Hook":
	case "Pipeline Hook":
	case "Build Hook":
		log.Warningf("Unhandled X-Gitlab-Event %s", eventType.Type)
		return nil
		break
	default:
		var e errorResponse
		e.Body.Message = "Invalid request"
		c.JSON(http.StatusBadRequest, e.Body)
		return nil
	}

	return nil
}
