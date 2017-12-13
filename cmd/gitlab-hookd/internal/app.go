package internal

import (
	"fmt"
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot"
)

var rabbitmqPublisher *bot.EventPublisher

// AppName application name
var AppName = "gitlab-hook"

// AppVersion application version (from git tag)
var AppVersion = "[unk]"

// AppBuildDate application build date
var AppBuildDate = "[unk]"

// StartApp initiate components
// Should be called from main function
func StartApp(configFile string) {
	initLogger()

	log.Infof("Starting %s version %s.", AppName, AppVersion)
	log.Infof("Build date: %s.", AppBuildDate)

	loadConfiguration(configFile)

	verifyPublisher()

	// Bind main thread to HTTP service
	e := echo.New()
	e.POST("/v1/gitlab/event", v1ApiGitlabEvent)

	httpListeningAddress := fmt.Sprintf(":%d", gconfig.HTTP.Port)

	e.Logger.Error(e.Start(httpListeningAddress))

	log.Infof("Exiting %s", AppName)
}
