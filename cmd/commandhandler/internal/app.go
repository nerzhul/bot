package internal

// AppName application name
var AppName = "commandhandler"

// AppVersion application version (from git tag)
var AppVersion = "[unk]"

// AppBuildDate application build date
var AppBuildDate = "[unk]"

var router *commandRouter

// StartApp initiate components
// Should be called from main function
func StartApp(configFile string) {
	initLogger()

	log.Infof("Starting %s version %s.", AppName, AppVersion)
	log.Infof("Build date: %s.", AppBuildDate)

	loadConfiguration(configFile)

	verifyPublisher()
	verifyConsumer()

	runProcessor()

	log.Infof("Exiting %s", AppName)
}
