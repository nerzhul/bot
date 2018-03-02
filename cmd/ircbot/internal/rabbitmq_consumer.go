package internal

import (
	"encoding/json"
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

	if ircConn == nil {
		msg.Nack(false, true)
		return
	}

	if !gconfig.isAllowedToUseCommand(command.User) {
		log.Errorf("User '%s' is not allowed to use IRC bot commands. Dropping.", command.User)
		msg.Ack(true)
		return
	}

	if len(command.Command) == 0 {
		log.Errorf("Ignore empty command received from user '%s'", command.User)
		msg.Ack(true)
		return
	}

	log.Debugf("Received command to handle '%s' from user '%s'", command.User)
	commandSpl := strings.Split(command.Command, " ")
	switch commandSpl[0] {
	case "join":
		if len(commandSpl) != 2 {
			log.Errorf("Command '%s' sent from user '%s' is malformed. 1 argument expected.",
				commandSpl[0], command.User)
			break
		}
		ircConn.Join(commandSpl[1])
		break
	case "leave":
		if len(commandSpl) != 2 {
			log.Errorf("Command '%s' sent from user '%s' is malformed. 1 argument expected.",
				commandSpl[0], command.User)
			break
		}
		ircConn.Part(commandSpl[1])
		break
	case "list":
		if len(commandSpl) != 1 {
			log.Errorf("Command '%s' sent from user '%s' is malformed. 0 argument expected.",
				commandSpl[0], command.User)
			break
		}
		// TODO: Not implemented
	default:
		log.Warningf("Ignore invalid command '%s' received from user '%s'", commandSpl[0], command.User)
		break
	}
	msg.Ack(true)
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

func (rc *rabbitmqClient) verifyConsumer() bool {
	if rc.consumer == nil {
		rc.consumer = rabbitmq.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rc.consumer.Init() {
			rc.consumer = nil
			return false
		}

		for _, consumerName := range []string{"ircbot", "chat"} {
			consumerCfg := gconfig.RabbitMQ.GetConsumer(consumerName)
			if consumerCfg == nil {
				log.Fatalf("RabbitMQ consumer configuration '%s' not found, aborting.", consumerName)
			}

			if !rc.consumer.DeclareQueue(consumerCfg.Queue) {
				rc.consumer = nil
				return false
			}

			if !rc.consumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
				rc.consumer = nil
				return false
			}

			if !rc.consumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, consumeResponses, false) {
				rc.consumer = nil
				return false
			}
		}
	}

	return rc.consumer != nil
}
