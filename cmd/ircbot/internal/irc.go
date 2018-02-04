package internal

import (
	"crypto/tls"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot"
	"strings"
	"time"
)

var ircConn *irc.Conn
var ircDisconnected chan bool

func joinConfiguredChannels(conn *irc.Conn) {
	for _, channel := range gconfig.IRC.Channels {
		if len(channel.Password) > 0 {
			log.Infof("Try joining channel %s (with password) on %s", channel.Name, conn.Config().Server)
			conn.Join(channel.Name, channel.Password)
		} else {
			log.Infof("Try joining channel %s on %s", channel.Name, conn.Config().Server)
			conn.Join(channel.Name)
		}
	}
}

func onIRCConnected(conn *irc.Conn, line *irc.Line) {
	log.Infof("Connected to IRC on %s", conn.Config().Server)
	// If we have a password, join later in the process
	if len(gconfig.IRC.Password) != 0 {
		return
	}

	joinConfiguredChannels(conn)
}

func onIRCDisconnected(conn *irc.Conn, line *irc.Line) {
	log.Infof("Disconnected from IRC on %s", conn.Config().Server)
	ircDisconnected <- true
}

func onIRCKick(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	log.Infof("Kicked from channel %s by %s, rejoining", line.Args[0], line.Nick)
	conn.Join(line.Args[0])
}

func onIRCJoin(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		return
	}

	if line.Nick == conn.Me().Nick {
		log.Infof("Channel %s joined on %s", line.Args[0], conn.Config().Server)
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

	if !verifyPublisher() {
		log.Error("Failed to verify publisher, no message sent to broker")
		return
	}

	if !verifyConsumer() {
		log.Error("Failed to verify consumer, no message sent to broker")
		return
	}

	// Publish chat event to handler
	rabbitmqPublisher.Publish(
		&bot.IRCChatEvent{
			Type:    "privmsg",
			Message: text,
			Channel: channel,
			User:    line.Nick,
		},
		"irc-chat",
		&bot.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  1800000,
			RoutingKey:    "irc-chat",
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

	ce := bot.CommandEvent{
		Command: text[1:],
		Channel: channel,
		User:    line.Nick,
	}

	log.Infof("User %s sent command on channel %s: %s", ce.User, ce.Channel, ce.Command)

	consumerCfg := gconfig.RabbitMQ.GetConsumer("ircbot")
	if consumerCfg == nil {
		log.Fatalf("RabbitMQ consumer configuration 'ircbot' not found, aborting.")
	}

	rabbitmqPublisher.Publish(
		&ce,
		"command",
		&bot.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ReplyTo:       consumerCfg.RoutingKey,
			ExpirationMs:  300000,
		},
	)
}

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

	if !verifyPublisher() {
		log.Error("Failed to verify publisher, no notice sent to broker")
		return
	}

	if !verifyConsumer() {
		log.Error("Failed to verify consumer, no notice sent to broker")
		return
	}

	channel := line.Args[0]

	if channel == conn.Me().Nick {
		channel = line.Nick
	}

	// Don't send global channel messages to broker
	if channel == "*" {
		return
	}

	// Publish chat event to handler
	rabbitmqPublisher.Publish(
		&bot.IRCChatEvent{
			Type:    "notice",
			Message: text,
			Channel: channel,
			User:    line.Nick,
		},
		"irc-chat",
		&bot.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ExpirationMs:  1800000,
			RoutingKey:    "irc-chat",
		},
	)
}

func onIRCError(conn *irc.Conn, line *irc.Line) {
	if len(line.Args) == 0 {
		log.Warningf("Received error %s", line.Text())
		return
	}

	text := line.Text()
	log.Warningf("Received error %s and args %v", text, line.Args)
}

func runIRCClient() {
	for {
		cfg := irc.NewConfig(gconfig.IRC.Name)
		cfg.SSL = true
		cfg.SSLConfig = &tls.Config{ServerName: gconfig.IRC.Server}
		cfg.Server = fmt.Sprintf("%s:%d", gconfig.IRC.Server, gconfig.IRC.Port)
		cfg.Me.Ident = gconfig.IRC.Name
		cfg.Me.Name = "For Ironforge"
		cfg.NewNick = func(n string) string { return n + "^" }
		ircConn = irc.Client(cfg)

		ircDisconnected = make(chan bool)

		ircConn.HandleFunc(irc.CONNECTED, onIRCConnected)
		ircConn.HandleFunc(irc.DISCONNECTED, onIRCDisconnected)
		ircConn.HandleFunc(irc.KICK, onIRCKick)
		ircConn.HandleFunc(irc.JOIN, onIRCJoin)
		ircConn.HandleFunc(irc.PRIVMSG, onIRCPrivMsg)
		ircConn.HandleFunc(irc.NOTICE, onIRCNotice)
		ircConn.HandleFunc(irc.ERROR, onIRCError)

		if err := ircConn.Connect(); err != nil {
			log.Errorf("Connection error: %s\n", err.Error())
			return
		}
		<-ircDisconnected
		ircConn = nil

		log.Errorf("Connection to IRC lost, retrying in 30sec")
		time.Sleep(time.Second * 30)
	}
}
