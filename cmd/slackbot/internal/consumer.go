package internal

import (
	"encoding/json"
	"github.com/nlopes/slack"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
)

var rabbitmqConsumer *bot.EventConsumer
var slackMsgID = 0

func consumeResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		if d.Type == "tweet" {
			consumeTwitterResponse(&d)
		} else {
			response := bot.CommandResponse{}
			err := json.Unmarshal(d.Body, &response)
			if err != nil {
				log.Errorf("Failed to decode command response : %v", err)
			}

			// Send message on slack
			slackMsgID++
			slackRTM.SendMessage(&slack.OutgoingMessage{
				ID:      slackMsgID,
				Type:    "message",
				Channel: response.Channel,
				Text:    response.Message,
			})

			d.Ack(false)
		}
	}
}

func consumeTwitterResponse(msg *amqp.Delivery) {
	tweet := bot.TweetMessage{}
	err := json.Unmarshal(msg.Body, &tweet)
	if err != nil {
		log.Errorf("Failed to decode tweet : %v", err)
	}

	// Send message on slack
	//slackMsgID++
	//slackRTM.SendMessage(&slack.OutgoingMessage{
	//	ID:      slackMsgID,
	//	Type:    "message",
	//	Channel: response.Channel,
	//	Text:    response.Message,
	//})

	msg.Ack(false)
}

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = bot.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.DeclareQueue(gconfig.RabbitMQ.EventQueue + "/commands") {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.BindExchange(gconfig.RabbitMQ.EventQueue+"/commands",
			gconfig.RabbitMQ.EventExchange, gconfig.RabbitMQ.ConsumerRoutingKey) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.DeclareQueue(gconfig.RabbitMQ.EventQueue + "/twitter") {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.BindExchange(gconfig.RabbitMQ.EventQueue+"/twitter",
			gconfig.TwitterRabbitMQExchange, gconfig.RabbitMQ.ConsumerRoutingKey) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.Consume(gconfig.RabbitMQ.EventQueue+"/commands", consumeResponses, false) {
			rabbitmqConsumer = nil
			return false
		}

		if !rabbitmqConsumer.Consume(gconfig.RabbitMQ.EventQueue+"/twitter", consumeResponses, false) {
			rabbitmqConsumer = nil
			return false
		}
	}

	return rabbitmqConsumer != nil
}
