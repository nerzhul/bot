package internal

import (
	"fmt"
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/gitlab"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/rabbitmq"
)

// AppName application name
var AppName = "webhook"

// AppVersion application version (from git tag)
var AppVersion = "[unk]"

// AppBuildDate application build date
var AppBuildDate = "[unk]"

// StartApp initiate components
// Should be called from main function
func StartApp(configFile string) {
	common.InitLogger(AppName)

	common.Log.Infof("Starting %s version %s.", AppName, AppVersion)
	common.Log.Infof("Build date: %s.", AppBuildDate)

	common.LoadConfiguration(configFile)

	rabbitmq.VerifyPublisher()

	// Bind main thread to HTTP service
	e := echo.New()
	e.POST("/v1/gitlab/event", gitlab.V1ApiGitlabEvent)
	if common.GConfig.Mattermost.EnableHook {
		e.POST("/v1/mattermost/command", v1ApiMattermostCommand)
	}

	httpListeningAddress := fmt.Sprintf(":%d", common.GConfig.HTTP.Port)

	e.Logger.Error(e.Start(httpListeningAddress))

	common.Log.Infof("Exiting %s", AppName)
}
