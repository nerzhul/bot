package internal

import (
	"crypto/tls"
	"fmt"
	irc "github.com/fluffle/goirc/client"
	"time"
)

var ircConn *irc.Conn
var ircDisconnected chan bool

type ircClient struct {
}

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

func (i *ircClient) run() {
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
