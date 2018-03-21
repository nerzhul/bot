package internal

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"strings"
)

func consumeResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		if d.Type == "irc-chat" {
			consumeIRCChatMessage(&d)
		} else if d.Type == "irc-command" {
			consumeIRCCommand(&d)
		} else {
			consumeCommandResponse(&d)
		}
	}
}

func consumeIRCChatMessage(msg *amqp.Delivery) {
	chatEvent := rabbitmq.IRCChatEvent{}
	err := json.Unmarshal(msg.Body, &chatEvent)
	if err != nil {
		log.Errorf("Failed to decode chat event: %v", err)
		msg.Nack(false, false)
		return
	}

	if ircConn == nil {
		msg.Nack(false, true)
		return
	}

	log.Debugf("Received message to send on IRC channel '%s': %s", chatEvent.Channel, chatEvent.Message)
	for _, msg := range strings.Split(chatEvent.Message, "\n") {
		ircConn.Privmsg(chatEvent.Channel, msg)
	}

	msg.Ack(false)
}

func consumeIRCCommand(msg *amqp.Delivery) {
	command := rabbitmq.IRCCommand{}
	err := json.Unmarshal(msg.Body, &command)
	if err != nil {
		log.Errorf("Failed to decode chat event: %v", err)
		msg.Nack(false, false)
		return
	}

	if len(command.User) == 0 {
		log.Error("IRCCommand user field is empty")
		msg.Nack(false, false)
		return
	}

	if len(command.Channel) == 0 {
		log.Error("IRCCommand channel field is empty")
		msg.Nack(false, false)
		return
	}

	resp := &rabbitmq.CommandResponse{
		Channel:     command.Channel,
		User:        command.User,
		MessageType: "whisper",
	}

	if !gconfig.isAllowedToUseCommand(command.User) {
		resp.Message = "You are not allowed to interact with IRC client. This will be reported."
		sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
		log.Errorf("User '%s' is not allowed to use IRC bot commands. Dropping.", command.User)
		msg.Nack(false, false)
		return
	}

	if ircConn == nil {
		resp.Message = "Unable to handle command, not connected to IRC."
		sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
		log.Warning("Received an IRC command whereas we are not connected, ignoring.")
		msg.Nack(false, false)
		return
	}

	if len(command.Command) == 0 {
		resp.Message = "Invalid command, ignoring."
		sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
		log.Errorf("Ignore empty command received from user '%s'", command.User)
		msg.Nack(false, false)
		return
	}

	log.Debugf("Received command to handle '%s' from user '%s'", command.Command, command.User)
	switch command.Command {
	case "join":
		if len(command.Arg1) == 0 {
			resp.Message = "Invalid command, ignoring."
			sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
			log.Errorf("Command '%s' sent from user '%s' is malformed. 1 argument expected.",
				command.Command, command.User)
			break
		}
		ircConn.Join(command.Arg1, command.Arg2)

		resp.Message = fmt.Sprintf("Bot requested to join channel '%s'.", command.Arg1)
		if err := gIRCDB.SaveIRCChannelConfig(command.Arg1, command.Arg2); err != nil {
			resp.Message += " But we failed to save the join state. It will be temporary."
		}
		sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
		break
	case "leave":
		if len(command.Arg1) == 0 {
			resp.Message = "Invalid command, ignoring."
			sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
			log.Errorf("Command '%s' sent from user '%s' is malformed. 1 argument expected.",
				command.Command, command.User)
			break
		}
		ircConn.Part(command.Arg1)

		resp.Message = fmt.Sprintf("Bot left channel '%s'.", command.Arg1)
		sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
		break
	case "list":
		if len(command.Arg1) != 0 {
			resp.Message = "Invalid command, ignoring."
			sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
			log.Errorf("Command '%s' sent from user '%s' is malformed. 0 argument expected.",
				command.Command, command.User)
			break
		}
		// TODO: Not implemented
	default:
		resp.Message = "Invalid command, ignoring."
		sendIRCCommandResponse(resp, msg.CorrelationId, msg.ReplyTo)
		log.Warningf("Ignore invalid command '%s' received from user '%s'", command.Command, command.User)
		break
	}

	msg.Ack(true)
}

func sendIRCCommandResponse(resp *rabbitmq.CommandResponse, correlationID string, replyTo string) {
	asyncClient.Publisher.Publish(resp,
		"irccommand-answer",
		&rabbitmq.EventOptions{
			CorrelationID: correlationID,
			RoutingKey:    replyTo,
			ExpirationMs:  60 * 1000,
		},
	)
}

func consumeCommandResponse(msg *amqp.Delivery) {
	response := rabbitmq.CommandResponse{}
	err := json.Unmarshal(msg.Body, &response)
	if err != nil {
		log.Errorf("Failed to decode command response : %v", err)
		msg.Nack(false, false)
		return
	}

	if ircConn == nil {
		msg.Nack(false, true)
		return
	}

	for _, msg := range strings.Split(response.Message, "\n") {
		if response.MessageType == "notice" {
			ircConn.Notice(response.Channel, msg)
		} else {
			ircConn.Privmsg(response.Channel, msg)
		}
	}

	msg.Ack(false)
}
