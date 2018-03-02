package internal

import (
	irc "github.com/fluffle/goirc/client"
)

func onIRCError(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		log.Warningf("Received error %s", line.Text())
		return
	}

	text := line.Text()
	log.Warningf("Received error %s and args %v", text, line.Args)
}
