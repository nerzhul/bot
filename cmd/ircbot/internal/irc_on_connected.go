package internal

import (
	irc "github.com/fluffle/goirc/client"
)

func onIRCConnected(conn *irc.Conn, line *irc.Line) {
	log.Infof("Connected to IRC on %s", conn.Config().Server)
	// If we have a password, join later in the process
	if len(gconfig.IRC.Password) != 0 {
		return
	}

	joinConfiguredChannels(conn)
}
