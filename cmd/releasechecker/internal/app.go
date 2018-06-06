package internal

import (
	"gitlab.com/nerzhul/bot/utils"
	"os"
	"os/signal"
	"syscall"
)

// AppName application name
var AppName = "releasechecker"

// AppVersion application version (from git tag)
var AppVersion = "[unk]"

// AppBuildDate application build date
var AppBuildDate = "[unk]"

var gDB *rcDB

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

	gDB = &rcDB{
		config: &gconfig.DB,
	}
	if !gDB.init() {
		log.Fatal("Failed to initialize database connector, aborting.")
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP)

	go func() {
		for sig := range sigs {
			log.Infof("SIGHUP(%s) received, reloading configuration", sig)
			loadConfiguration(configFile)
		}
	}()

	if !checkGithubNewTags() {
		os.Exit(1)
	}

	log.Infof("%s has finished.", AppName)
}
