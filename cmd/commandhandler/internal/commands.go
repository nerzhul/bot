package internal

import "gitlab.com/nerzhul/bot"

type commandHandler struct {
}

var commandHandlers = map[string]commandHandler{
	"start-ci": {},
}

func handleCommand(event *bot.CommandEvent) bool {
	log.Infof("Receive command event from user %s", event.User)
	return true
}
