package internal

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"strings"
)

var rabbitmqConsumer *rabbitmq.EventConsumer

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
	}

	if ircConn == nil {
		msg.Nack(false, true)
	}

	log.Debugf("Received message to send on IRC channel '%s': %s", chatEvent.Channel, chatEvent.Message)
	for _, msg := range strings.Split(chatEvent.Message, "\n") {
		ircConn.Privmsg(chatEvent.Channel, msg)
	}

	msg.Ack(false)
}

func consumeIRCCommand(msg *amqp.Delivery) {

}

func consumeCommandResponse(msg *amqp.Delivery) {
	response := rabbitmq.CommandResponse{}
	err := json.Unmarshal(msg.Body, &response)
	if err != nil {
		log.Errorf("Failed to decode command response : %v", err)
		msg.Nack(false, false)
	}

	if ircConn == nil {
		msg.Nack(false, true)
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

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = rabbitmq.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		for _, consumerName := range []string{"ircbot", "chat"} {
			consumerCfg := gconfig.RabbitMQ.GetConsumer(consumerName)
			if consumerCfg == nil {
				log.Fatalf("RabbitMQ consumer configuration '%s' not found, aborting.", consumerName)
			}

			if !rabbitmqConsumer.DeclareQueue(consumerCfg.Queue) {
				rabbitmqConsumer = nil
				return false
			}

			if !rabbitmqConsumer.BindExchange(consumerCfg.Queue, consumerCfg.Exchange, consumerCfg.RoutingKey) {
				rabbitmqConsumer = nil
				return false
			}

			if !rabbitmqConsumer.Consume(consumerCfg.Queue, consumerCfg.ConsumerID, consumeResponses, false) {
				rabbitmqConsumer = nil
				return false
			}
		}
	}

	return rabbitmqConsumer != nil
}
