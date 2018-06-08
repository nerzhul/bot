package internal

import (
	"encoding/json"
	//"fmt"
	"fmt"
	"github.com/mattermost/mattermost-server/model"
	"github.com/streadway/amqp"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"strings"
	"time"
)

var rabbitmqConsumer *rabbitmq.EventConsumer

func consumeResponses(msg *amqp.Delivery) {
	if msg.Type == "irc-chat" {
		consumeIRCResponse(msg)
	} else if msg.Type == "tweet" {
		consumeTwitterResponse(msg)
	} else if msg.Type == "announcement" {
		consumeAnnouncementMessage(msg)
	} else {
		consumeCommandResponse(msg)
	}
}

func consumeAnnouncementMessage(msg *amqp.Delivery) {
	if !mClient.isMattermostUp() {
		log.Warningf("Mattermost client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	announceMsg := rabbitmq.AnnouncementMessage{}
	err := json.Unmarshal(msg.Body, &announceMsg)
	if err != nil {
		log.Errorf("Failed to decode announcement message : %v", err)
	}

	postMessageToMattermostViaWebhook(gconfig.Mattermost.ReleaseAnnouncementsChannel,
		fmt.Sprintf("**%s** [%s](%s) has been released !",
			announceMsg.What, announceMsg.Message, announceMsg.URL), "ReleaseChecker", msg)
}

func consumeCommandResponse(msg *amqp.Delivery) {
	if !mClient.isMattermostUp() {
		log.Warningf("Mattermost client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	response := rabbitmq.CommandResponse{}
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
		log.Errorf("Failed to send a message to '%s' channel: %s", response.Channel, resp.Error.Message)
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

	ircChatEvent := rabbitmq.IRCChatEvent{}
	err := json.Unmarshal(msg.Body, &ircChatEvent)
	if err != nil {
		log.Errorf("Failed to decode tweet : %v", err)
		msg.Nack(false, false)
		return
	}

	if ircChatEvent.Channel == "*" || ircChatEvent.Channel == "$$*" {
		msg.Ack(false)
		return
	}

	channelDisplayName := fmt.Sprintf("irc-%s", ircChatEvent.Channel)
	channelName := strings.Replace(channelDisplayName, "#", "", -1)
	mClient.createChannelIfNeeded(channelName, channelDisplayName, model.CHANNEL_OPEN)

	chanInfo := mClient.getChannelInfo(channelName)
	if chanInfo == nil {
		log.Errorf("Unable to find mattermost channel %s", channelName)
		msg.Nack(false, true)
		return
	}

	if ircChatEvent.Type == "privmsg" || ircChatEvent.Type == "notice" {
		handleIRCChatEventMessage(&ircChatEvent, channelName, msg)
	} else if ircChatEvent.Type == "topic" {
		header := new(string)
		*header = ircChatEvent.Message

		_, response := mClient.client.PatchChannel(chanInfo.Id, &model.ChannelPatch{
			Header: header,
		})

		if response.Error != nil {
			log.Errorf("Failed to update topic. Error was: %s", response.Error.Message)
		}
	} else {
		log.Warningf("Ignore unknown irc chat event type '%s'", ircChatEvent.Type)
		msg.Nack(false, false)
	}

}

func handleIRCChatEventMessage(ircChatEvent *rabbitmq.IRCChatEvent, channelName string, msg *amqp.Delivery) {
	postMessageToMattermostViaWebhook(channelName, ircChatEvent.Message, ircChatEvent.User, msg)
}

func consumeTwitterResponse(msg *amqp.Delivery) {
	if !mClient.isMattermostUp() {
		log.Warningf("Mattermost client is not ready, waiting 1sec...")
		time.Sleep(time.Second * 1)
		msg.Nack(false, true)
		return
	}

	tweet := rabbitmq.TweetMessage{}
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
		rabbitmqConsumer = rabbitmq.NewEventConsumer(log, &gconfig.RabbitMQ)
		if !rabbitmqConsumer.Init() {
			rabbitmqConsumer = nil
			return false
		}

		for _, consumerName := range []string{"announcements", "commands", "irc", "twitter"} {
			consumerCfg := gconfig.RabbitMQ.GetConsumer(consumerName)
			if consumerCfg == nil {
				log.Fatalf("RabbitMQ consumer configuration '%s' not found, aborting.", consumerName)
			}

			if !rabbitmqConsumer.DeclareQueue(consumerCfg.Queue) {
				rabbitmqConsumer = nil
				return false
			}

			if !rabbitmqConsumer.DeclareExchange(consumerCfg.Exchange, consumerCfg.ExchangeType,
				consumerCfg.ExchangeDurable) {
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
