package common

import (
	"github.com/op/go-logging"
	"gitlab.com/nerzhul/bot/utils"
	"os"
)

// Log global logger
var Log = logging.MustGetLogger("webhook")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} - %{level:.5s} %{color:reset} %{message}`,
)

// InitLogger initialize logger
func InitLogger(name string) {
	if !utils.IsInDocker() {
		stderrLog := logging.NewLogBackend(os.Stderr, "", 0)
		syslogBackend, err := logging.NewSyslogBackend(name)
		if err != nil {
			Log.Error("Failed to setup logs syslog backend.")
		}
		logging.SetBackend(logging.NewBackendFormatter(stderrLog, format), syslogBackend)
	}
}
