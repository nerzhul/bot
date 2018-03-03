package internal

import (
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"strings"
)

func onIRCNotice(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	text := line.Text()

	if line.Nick == "NickServ" {
		if strings.Contains(text, "This nickname is registered") {
			log.Infof("Authentication request from NickServ on %s", conn.Config().Server)
			conn.Privmsg(line.Nick, fmt.Sprintf("IDENTIFY %s", gconfig.IRC.Password))
		} else if strings.Contains(text, "You are now identified for") {
			log.Infof("Authentication succeed on %s.", conn.Config().Server)
			joinConfiguredChannels(conn)
		} else if strings.Contains(text, "Invalid password for") {
			log.Infof("Authentication failed on %s, disconnecting.", conn.Config().Server)
			conn.Close()
		}

		return
	}

	if !asyncClient.VerifyPublisher() {
		log.Error("Failed to verify publisher, no notice sent to broker")
		return
	}

	if !asyncClient.VerifyConsumer() {
		log.Error("Failed to verify consumer, no notice sent to broker")
		return
	}

	channel := line.Args[0]

	if channel == conn.Me().Nick {
		channel = line.Nick
	}

	// Don't send global channel messages to broker
	if channel == "*" || channel == "$$*" {
		return
	}

	// Publish chat event to handler
	asyncClient.publishChatEvent(
		&rabbitmq.IRCChatEvent{
			Type:    "notice",
			Message: text,
			Channel: channel,
			User:    line.Nick,
		},
	)
}
