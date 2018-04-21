package internal

import (
	"github.com/op/go-logging"
	"gitlab.com/nerzhul/bot/utils"
	"os"
)

var log = logging.MustGetLogger(AppName)
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} - %{level:.5s} %{color:reset} %{message}`,
)

func initLogger() {
	if !utils.IsInDocker() {
		stderrLog := logging.NewLogBackend(os.Stderr, "", 0)
		syslogBackend, err := logging.NewSyslogBackend(AppName)
		if err != nil {
			log.Error("Failed to setup logs syslog backend.")
		}
		logging.SetBackend(logging.NewBackendFormatter(stderrLog, format), syslogBackend)
	}
}
