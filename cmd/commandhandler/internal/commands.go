package internal

import (
	"fmt"
	"gitlab.com/nerzhul/bot"
	"strings"
)

// commandHandler (args, user, channel)
// returns nil on failure else the response string
type commandHandler func(string, string, string) *string

type commandRouter struct {
	commandList     []string
	commandHandlers map[string]commandHandler
}

func (r *commandRouter) init() {
	log.Infof("Initialize command router...")

	r.commandHandlers = map[string]commandHandler{
		"help": r.handlerHelp,
	}

	r.commandList = []string{}
	for k := range r.commandHandlers {
		r.commandList = append(r.commandList, k)
	}

	log.Infof("Router init done (%d commands registered).", len(r.commandList))
}

func (r *commandRouter) handleCommand(event *bot.CommandEvent, correlationID string, replyTo string) bool {
	log.Infof("Receive command event from user %s", event.User)
	ecsplit := strings.SplitN(event.Command, " ", 2)
	if len(ecsplit) == 0 {
		log.Errorf("Failed to split command '%s', ignoring command.", event.Command)
		return true
	}

	command := ecsplit[0]
	commandArgs := ""
	if len(ecsplit) > 2 {
		log.Fatalf("SplitN command length > 2, aborting.")
	} else if len(ecsplit) == 2 {
		commandArgs = ecsplit[1]
	}

	log.Infof("Command %s (args: '%s') sent by %s on channel %s", command, commandArgs, event.User, event.Channel)

	resp := bot.CommandResponse{
		Channel: event.Channel,
		User:    event.User,
	}
	if val, ok := r.commandHandlers[command]; ok {
		// Execute command callback
		cmdResult := val(commandArgs, event.User, event.Channel)
		if cmdResult == nil {
			resp.Message = fmt.Sprintf("Failed to process command %s, verify server logs.", command)
		} else {
			resp.Message = *cmdResult
		}
	} else {
		if !verifyPublisher() {
			return false
		}

		resp.Message = "Invalid command. Call help command to know the available commands."
	}

	rabbitmqPublisher.Publish(&resp,
		"command-answer",
		&bot.EventOptions{
			CorrelationID: correlationID,
			RoutingKey:    replyTo,
			ExpirationMs:  60 * 1000,
		},
	)
	return true
}
