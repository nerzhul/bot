package internal

import (
	irc "github.com/fluffle/goirc/client"
)

func onIRCUser(conn *irc.Conn, line *irc.Line) {
	log.Debugf("%v", line)
}
