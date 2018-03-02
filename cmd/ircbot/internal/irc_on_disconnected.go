package internal

import (
	irc "github.com/fluffle/goirc/client"
)

func onIRCDisconnected(conn *irc.Conn, line *irc.Line) {
	log.Infof("Disconnected from IRC on %s", conn.Config().Server)
	ircDisconnected <- true
}
