package main

import (
	"github.com/pborman/getopt/v2"
	"gitlab.com/nerzhul/bot/cmd/releasechecker/internal"
)

var configFile = ""

func init() {
	getopt.FlagLong(&configFile, "config", 'c', "Configuration file")
}

func main() {
	getopt.Parse()
	internal.StartApp(configFile)
}
