package internal

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

func onIRCJoin(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	if conn.StateTracker() == nil {
		log.Fatal("IRC connection state tracker is nil")
	}

	if line.Nick == conn.Me().Nick {
		log.Infof("Channel %s joined on %s", line.Args[0], conn.Config().Server)

		if !asyncClient.VerifyPublisher() {
			log.Error("Failed to verify publisher, no message sent to broker")
			return
		}

		// Publish topic event to queue
		asyncClient.publishChatEvent(
			&rabbitmq.IRCChatEvent{
				Type:    "topic",
				Message: conn.StateTracker().GetChannel(line.Args[0]).Topic,
				Channel: line.Args[0],
				User:    line.Nick,
			},
		)

		// publish join event to queue
		asyncClient.publishChatEvent(
			&rabbitmq.IRCChatEvent{
				Type:    "notice",
				Message: fmt.Sprintf("Channel '%s' joined by the bot", line.Args[0]),
				Channel: line.Args[0],
				User:    line.Nick,
			},
		)
	}

	channelCfg := gconfig.getIRCChannelConfig(line.Args[0])
	if channelCfg == nil || !channelCfg.Hello {
		return
	}

	if line.Nick == conn.Me().Nick {
		conn.Privmsg(line.Args[0], fmt.Sprintf("Hello %s!", line.Args[0]))
	} else {
		conn.Privmsg(line.Args[0], fmt.Sprintf("Hello %s!", line.Nick))
	}
}
