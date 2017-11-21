package internal

import (
	"github.com/op/go-logging"
)

var log = logging.MustGetLogger("uc-achievements")
var format = logging.MustStringFormatter(
	`%{color}%{time:15:04:05.000} %{shortfunc} - %{level:.5s} %{color:reset} %{message}`,
)

func GetLogger() *logging.Logger {
	return log
}
