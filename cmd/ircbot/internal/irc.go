package internal

import (
	"crypto/tls"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"github.com/satori/go.uuid"
	"strings"
	"time"
)

var ircConn *irc.Conn
var ircDisconnected chan bool

func onIRCConnected(conn *irc.Conn, line *irc.Line) {
	log.Infof("Connected to IRC on %s", conn.Config().Server)
	for _, channel := range gconfig.IRC.Channels {
		if len(channel.Password) > 0 {
			conn.Join(channel.Name, channel.Password)
		} else {
			conn.Join(channel.Name)
		}
	}
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
	if len(text) < 2 || text[0] != '!' {
		return
	}

	channel := line.Args[0]

	// If it's a private message, channel is user
	if channel == conn.Me().Nick {
		channel = line.Nick
	} else {
		// We are on a channel, verify we answer commands
		channelCfg := gconfig.getIRCChannelConfig(channel)
		if channelCfg == nil || !channelCfg.AnswerCommands {
			return
		}
	}

	ce := commandEvent{
		Command: text[1:],
		Channel: channel,
		User:    line.Nick,
	}

	log.Infof("User %s sent command on channel %s: %s", ce.User, ce.Channel, ce.Command)

	if !verifyPublisher() {
		log.Error("Failed to verify publisher, no command sent to broker")
		return
	}

	if !verifyConsumer() {
		log.Error("Failed to verify consumer, no command sent to broker")
		return
	}

	rabbitmqPublisher.Publish(
		&ce,
		"command",
		uuid.NewV4().String(),
		gconfig.RabbitMQ.ConsumerRoutingKey,
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
			log.Infof("Authentication succeed on %s", conn.Config().Server)
		} else if strings.Contains(text, "Invalid password for") {
			log.Infof("Authentication failed on %s", conn.Config().Server)
		}
	}
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
