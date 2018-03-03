package internal

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

func onIRCPart(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	if line.Nick == conn.Me().Nick {
		log.Infof("Channel %s left on %s", line.Args[0], conn.Config().Server)

		if !asyncClient.VerifyPublisher() {
			log.Error("Failed to verify publisher, no message sent to broker")
			return
		}

		asyncClient.publishChatEvent(
			&rabbitmq.IRCChatEvent{
				Type:    "notice",
				Message: fmt.Sprintf("Channel '%s' left by the bot", line.Args[0]),
				Channel: line.Args[0],
				User:    line.Nick,
			},
		)
	}
}
