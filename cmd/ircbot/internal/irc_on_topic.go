package internal

import (
	irc "github.com/fluffle/goirc/client"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

func onIRCTopic(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) < 2 {
		log.Warningf("Ignore invalid TOPIC command. (%d < 2 arguments)", len(line.Args))
		return
	}

	log.Infof("Received topic '%s' for channel '%s'. Text: %s", line.Args[1], line.Args[0], line.Text())

	if !asyncClient.VerifyPublisher() {
		log.Error("Failed to verify publisher, no message sent to broker")
		return
	}

	if !asyncClient.VerifyConsumer() {
		log.Error("Failed to verify consumer, no message sent to broker")
		return
	}

	channel := line.Args[0]

	// Publish chat event to handler
	asyncClient.publishChatEvent(
		&rabbitmq.IRCChatEvent{
			Type:    "topic",
			Message: line.Args[1],
			Channel: channel,
			User:    line.Nick,
		},
	)
}
