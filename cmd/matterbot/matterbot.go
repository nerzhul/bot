// main package
//     Schemes: http, https
//     Host: localhost
//     BasePath: /
//     Version: 1.0
//     License: BSD
//     Contact: Support<support@unix-experience.fr> https://www.unix-experience.fr
//
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
// swagger:meta
package main

import (
	"github.com/pborman/getopt/v2"
	"gitlab.com/nerzhul/bot/cmd/matterbot/internal"
)

var configFile = ""

func init() {
	getopt.FlagLong(&configFile, "config", 'c', "Configuration file")
}

func main() {
	getopt.Parse()
	internal.StartApp(configFile)
}
