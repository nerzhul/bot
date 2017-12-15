package internal

import (
	"encoding/json"
	"fmt"
	"github.com/nlopes/slack"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
	"time"
)

var rabbitmqConsumer *bot.EventConsumer
var slackMsgID = 0

func consumeResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		if d.Type == "tweet" {
			consumeTwitterResponse(&d)
		} else {
			consumeCommandResponse(&d)
		}
	}
}

func consumeCommandResponse(msg *amqp.Delivery) {
	if slackRTM == nil || slackAPI == nil {
		log.Warningf("Slack client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	response := bot.CommandResponse{}
	err := json.Unmarshal(msg.Body, &response)
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

	msg.Ack(false)
}

func consumeTwitterResponse(msg *amqp.Delivery) {
	if slackRTM == nil || slackAPI == nil {
		log.Warningf("Slack client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	tweet := bot.TweetMessage{}
	err := json.Unmarshal(msg.Body, &tweet)
	if err != nil {
		log.Errorf("Failed to decode tweet : %v", err)
		msg.Nack(false, false)
		return
	}

	slackTweet := fmt.Sprintf("<https://twitter.com/%s|%s> @%s \n%s\n",
		tweet.UserScreenName, tweet.Username, tweet.UserScreenName, tweet.Message)

	params := slack.PostMessageParameters{}
	attachment := slack.Attachment{
		Text:  "Actions",
		Color: "#1da1f2",
		Actions: []slack.AttachmentAction{
			{
				Name:  "retweet",
				Type:  "button",
				Text:  "Retweet",
				Value: "retweet",
			},
		},
	}
	params.Attachments = []slack.Attachment{attachment}

	// Send message on slack
	slackAPI.PostMessage(gconfig.Slack.TwitterChannel, slackTweet, params)

	msg.Ack(false)
}

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = bot.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		for _, consumerName := range []string{"commands", "twitter"} {
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
