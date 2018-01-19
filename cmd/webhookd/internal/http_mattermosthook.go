package internal

import (
	"github.com/labstack/echo"
	"net/http"
)

// swagger:route POST /v1/mattermost/commands gitlab pushGitlabEvent
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
func v1ApiMattermostCommand(c echo.Context) error {
	return c.JSON(http.StatusBadRequest, nil)
}
