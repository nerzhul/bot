package internal

import (
	"os"
	"os/signal"
	"syscall"
)

// AppName application name
var AppName = "ircbot"

// AppVersion application version (from git tag)
var AppVersion = "[unk]"

// AppBuildDate application build date
var AppBuildDate = "[unk]"

var asyncClient *rabbitmqClient

// StartApp initiate components
// Should be called from main function
func StartApp(configFile string) {
	initLogger()

	log.Infof("Starting %s version %s.", AppName, AppVersion)
	log.Infof("Build date: %s.", AppBuildDate)

	loadConfiguration(configFile)

	asyncClient := newRabbitMQClient()

	asyncClient.VerifyPublisher()
	asyncClient.verifyConsumer()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for sig := range sigs {
			log.Infof("SIGHUP(%s) received, reloading configuration", sig)
			loadConfiguration(configFile)
		}
	}()

	irc := ircClient{}
	irc.run()

	log.Infof("Exiting %s", AppName)
}
