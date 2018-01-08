package internal

import (
	"gitlab.com/nerzhul/bot"
)

type commandHandler struct {
}

type commandRouter struct {
	commandList     []string
	commandHandlers map[string]commandHandler
}

var router *commandRouter

func (r *commandRouter) init() {
	log.Infof("Initialize command router...")

	r.commandHandlers = map[string]commandHandler{
		"start-builder": {},
	}

	r.commandList = []string{}
	for k := range r.commandHandlers {
		r.commandList = append(r.commandList, k)
	}

	log.Infof("Router init done (%d commands registered).", len(r.commandList))
}

func (r *commandRouter) handleCommand(event *bot.CommandEvent, correlationID string, replyTo string) bool {
	log.Infof("Receive command event from user %s", event.User)
	if val, ok := r.commandHandlers[event.Command]; ok {
		log.Infof("val %v", val)
	} else {
		if !verifyPublisher() {
			return false
		}

		errorResponse := bot.CommandResponse{
			Channel: event.Channel,
			User:    event.User,
			Message: "Invalid command",
		}

		rabbitmqPublisher.Publish(&errorResponse,
			"command-answer",
			&bot.EventOptions{
				CorrelationID: correlationID,
				RoutingKey:    replyTo,
				ExpirationMs:  60 * 1000,
			})
	}
	return true
}
