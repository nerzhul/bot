package internal

import "gitlab.com/nerzhul/gitlab-hook"

var rabbitmqPublisher *bot.EventPublisher

// AppName application name
var AppName = "slackbot"

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

	runSlackClient()

	log.Infof("Exiting %s", AppName)
}
