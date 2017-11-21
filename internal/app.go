package internal

import (
	"fmt"
	"github.com/labstack/echo"
)

var rabbitmqPublisher *gitlabEventPublisher

func StartApp(configFile string) {
	loadConfiguration(configFile)

	rabbitmqPublisher := NewGitlabEventPublisher()
	if !rabbitmqPublisher.Init() {
		rabbitmqPublisher = nil
	}

	// Bind main thread to HTTP service
	e := echo.New()
	e.POST("/v1/gitlab/event", v1ApiGitlabEvent)

	httpListeningAddress := fmt.Sprintf(":%d", gconfig.Http.Port)

	e.Logger.Error(e.Start(httpListeningAddress))
}
