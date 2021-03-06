package internal

import (
	"gitlab.com/nerzhul/bot/utils"
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

var gIRCDB *ircDB

// StartApp initiate components
// Should be called from main function
func StartApp(configFile string) {
	initLogger()

	log.Infof("Starting %s version %s.", AppName, AppVersion)
	log.Infof("Build date: %s.", AppBuildDate)
	if utils.IsInDocker() {
		log.Infof("Application is running in a Docker container.")
	}

	loadConfiguration(configFile)

	gIRCDB = &ircDB{
		config: &gconfig.DB,
	}
	if !gIRCDB.init() {
		log.Fatal("Failed to initialize database connector, aborting.")
	}

	gconfig.loadDatabaseConfigurations()

	asyncClient = newRabbitMQClient()
	asyncClient.AddConsumerName("ircbot")
	asyncClient.AddConsumerName("chat")

	asyncClient.VerifyPublisher()
	asyncClient.VerifyConsumer()

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
