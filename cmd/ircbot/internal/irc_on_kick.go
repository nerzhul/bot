package internal

import (
	irc "github.com/fluffle/goirc/client"
)

func onIRCKick(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	log.Infof("Kicked from channel %s by %s, rejoining", line.Args[0], line.Nick)
	conn.Join(line.Args[0])
}
