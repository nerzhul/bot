package internal

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"github.com/mattermost/mattermost-server/model"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot"
	"strings"
	"time"
)

var rabbitmqConsumer *bot.EventConsumer

func consumeResponses(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		if d.Type == "irc-chat" {
			consumeIRCResponse(&d)
		} else if d.Type == "tweet" {
			consumeTwitterResponse(&d)
		} else {
			consumeCommandResponse(&d)
		}
	}
}

func consumeCommandResponse(msg *amqp.Delivery) {
	if !mClient.isMattermostUp() {
		log.Warningf("Mattermost client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	response := bot.CommandResponse{}
	err := json.Unmarshal(msg.Body, &response)
	if err != nil {
		log.Errorf("Failed to decode command response : %v", err)
	}

	// Send message on mattermost
	post := &model.Post{
		ChannelId: response.Channel,
		Message:   response.Message,
	}

	if _, resp := mClient.client.CreatePost(post); resp.Error != nil {
		log.Errorf("Failed to send a message to '%s' channel.", response.Channel)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func consumeIRCResponse(msg *amqp.Delivery) {
	if !mClient.isMattermostUp() {
		log.Warningf("Mattermost client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	ircChatEvent := bot.IRCChatEvent{}
	err := json.Unmarshal(msg.Body, &ircChatEvent)
	if err != nil {
		log.Errorf("Failed to decode tweet : %v", err)
		msg.Nack(false, false)
		return
	}

	log.Debugf("Received IRC event %v", ircChatEvent)

	channelName := strings.Replace(fmt.Sprintf("irc-%s", ircChatEvent.Channel), "#", "", -1)
	mClient.createChannelIfNeeded(channelName, model.CHANNEL_OPEN)

	chanInfo := mClient.getChannelInfo(channelName)
	if chanInfo == nil {
		log.Errorf("Unable to find mattermost channel %s", channelName)
		msg.Nack(false, true)
		return
	}

	post := &model.Post{
		ChannelId: chanInfo.Id,
		Message:   fmt.Sprintf("Message from %s:\n%s", ircChatEvent.User, ircChatEvent.Message),
		Props: model.StringInterface{
			"username": ircChatEvent.User,
		},
	}

	if _, resp := mClient.client.CreatePost(post); resp.Error != nil {
		log.Errorf("Failed to send a message to '%s' channel.", channelName)
		msg.Nack(false, true)
		return
	}

	msg.Ack(false)
}

func consumeTwitterResponse(msg *amqp.Delivery) {
	if !mClient.isMattermostUp() {
		log.Warningf("Mattermost client is not ready, waiting 1sec...")
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

	//slackTweet := fmt.Sprintf("<https://twitter.com/%s|%s> @%s \n%s\n",
	//	tweet.UserScreenName, tweet.Username, tweet.UserScreenName, tweet.Message)
	//
	//params := slack.PostMessageParameters{}
	//attachment := slack.Attachment{
	//	Text:  "Actions",
	//	Color: "#1da1f2",
	//	Actions: []slack.AttachmentAction{
	//		{
	//			Name:  "retweet",
	//			Type:  "button",
	//			Text:  "Retweet",
	//			Value: "retweet",
	//		},
	//	},
	//}
	//params.Attachments = []slack.Attachment{attachment}

	// Send message on mattermost
	//slackAPI.PostMessage(gconfig.Slack.TwitterChannel, slackTweet, params)

	msg.Ack(false)
}

func verifyConsumer() bool {
	if rabbitmqConsumer == nil {
		rabbitmqConsumer = bot.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		for _, consumerName := range []string{"commands", "irc", "twitter"} {
			consumerCfg := gconfig.RabbitMQ.GetConsumer(consumerName)
			if consumerCfg == nil {
				log.Fatalf("RabbitMQ consumer configuration '%s' not found, aborting.", consumerName)
			}

			if !rabbitmqConsumer.DeclareQueue(consumerCfg.Queue) {
				rabbitmqConsumer = nil
				return false
			}

			if !rabbitmqConsumer.DeclareExchange(consumerCfg.Exchange, consumerCfg.ExchangeDurable) {
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
