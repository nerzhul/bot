package internal

import (
	irc "github.com/fluffle/goirc/client"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

func onIRCPrivMsg(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	text := line.Text()

	if len(text) < 2 {
		return
	}

	channel := line.Args[0]

	if channel == conn.Me().Nick {
		channel = line.Nick
	}

	if !asyncClient.verifyPublisher() {
		log.Error("Failed to verify publisher, no message sent to broker")
		return
	}

	if !asyncClient.verifyConsumer() {
		log.Error("Failed to verify consumer, no message sent to broker")
		return
	}

	// Publish chat event to handler
	asyncClient.publishChatEvent(
		&rabbitmq.IRCChatEvent{
			Type:    "privmsg",
			Message: text,
			Channel: channel,
			User:    line.Nick,
		},
	)

	// Don't send non commands to commandhandler
	if text[0] != '!' {
		return
	}

	if channel != conn.Me().Nick {
		// We are on a channel, verify if we answer to commands
		channelCfg := gconfig.getIRCChannelConfig(channel)
		if channelCfg == nil || !channelCfg.AnswerCommands {
			return
		}
	}

	ce := rabbitmq.CommandEvent{
		Command: text[1:],
		Channel: channel,
		User:    line.Nick,
	}

	log.Infof("User %s sent command on channel %s: %s", ce.User, ce.Channel, ce.Command)

	consumerCfg := gconfig.RabbitMQ.GetConsumer("ircbot")
	if consumerCfg == nil {
		log.Fatalf("RabbitMQ consumer configuration 'ircbot' not found, aborting.")
	}

	asyncClient.publishCommand(&ce, consumerCfg.RoutingKey)
}
