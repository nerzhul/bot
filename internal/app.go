package internal

import (
	"fmt"
	"github.com/labstack/echo"
)

func StartApp(configFile string) {
	loadConfiguration(configFile)

	// Bind main thread to HTTP service
	e := echo.New()
	e.POST("/v1/gitlab/event", v1ApiGitlabEvent)

	httpListeningAddress := fmt.Sprintf(":%d", gconfig.Http.Port)

	e.Logger.Error(e.Start(httpListeningAddress))
}
