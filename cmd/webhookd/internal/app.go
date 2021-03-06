package internal

import (
	"fmt"
	"github.com/labstack/echo"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/gitlab"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/mattermost"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/rabbitmq"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/slack"
	"gitlab.com/nerzhul/bot/utils"
	"os"
	"os/signal"
	"syscall"
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
	if utils.IsInDocker() {
		common.Log.Infof("Application is running in a Docker container.")
	}

	common.LoadConfiguration(configFile)

	rabbitmq.AsyncClient = rabbitmq.NewRabbitMQClient()
	rabbitmq.AsyncClient.AddConsumerName("webhook")
	rabbitmq.AsyncClient.VerifyPublisher()
	rabbitmq.AsyncClient.VerifyConsumer()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for sig := range sigs {
			common.Log.Infof("SIGHUP(%s) received, reloading configuration", sig)
			common.LoadConfiguration(configFile)
		}
	}()

	// Bind main thread to HTTP service
	e := echo.New()
	common.Log.Info("Binding URL /v1/gitlab/event")
	e.POST("/v1/gitlab/event", gitlab.V1ApiGitlabEvent)

	if common.GConfig.Mattermost.EnableHook {
		common.Log.Info("Binding URL /v1/mattermost/command")
		e.POST("/v1/mattermost/command", mattermost.V1ApiMattermostCommand)
	}

	if common.GConfig.Slack.EnableHook {
		common.Log.Info("Binding URL /v1/slack/command")
		e.POST("/v1/slack/command", slack.V1ApiSlackCommand)
	}

	httpListeningAddress := fmt.Sprintf(":%d", common.GConfig.HTTP.Port)

	e.Logger.Error(e.Start(httpListeningAddress))

	common.Log.Infof("Exiting %s", AppName)
}
