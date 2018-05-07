package internal

import (
	"encoding/json"
	"fmt"
	"github.com/mattermost/mattermost-server/model"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"strings"
	"time"
)

type mattermostClient struct {
	client *model.Client4
	user   *model.User
	team   *model.Team
	//channels map[string]*model.Channel
}

type mattermostWebhookEvent struct {
	Text     string `json:"text"`
	Username string `json:"username"`
	Channel  string `json:"channel"`
}

// toJSON convert achievement event to byte string
// Returns non nil error when marshaling failed
func (ae *mattermostWebhookEvent) toJSON() ([]byte, error) {
	jsonStr, err := json.Marshal(ae)
	if err != nil {
		return nil, err
	}

	return jsonStr, nil
}

var mClient mattermostClient

func runMattermostClient() {
	for {
		mClient.init()
		if mClient.isMattermostUp() {
			mClient.createChannelIfNeeded(gconfig.Mattermost.TwitterChannel,
				gconfig.Mattermost.TwitterChannel, model.CHANNEL_PRIVATE)
		}

		mClient.run()

		mClient.deinit()
		log.Warning("Connection to mattermost lost, retrying in 15s...")
		time.Sleep(time.Second * 15)
	}
}

func (m *mattermostClient) init() {
	m.client = model.NewAPIv4Client(gconfig.Mattermost.URL)
	//m.channels = make(map[string]*model.Channel)
}

func (m *mattermostClient) deinit() {
	m.client = nil
	m.user = nil
	m.team = nil
	//m.channels = nil
}

func (m *mattermostClient) login() bool {
	var resp *model.Response
	if m.user, resp = m.client.Login(gconfig.Mattermost.Email, gconfig.Mattermost.Password); resp.Error != nil {
		log.Error("There was a problem logging into the Mattermost server. Ensure login is correct.")
		return false
	}

	return true
}

func (m *mattermostClient) run() bool {
	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(gconfig.Mattermost.WsURL, m.client.AuthToken)
	if err != nil {
		log.Error("Failed to connect to Mattermost websocket.")
		return false
	}

	webSocketClient.Listen()

	for {
		select {
		case resp := <-webSocketClient.EventChannel:
			if !m.handleWebSocketEvent(resp) {
				goto end
			}
		}
	}

end:
	if webSocketClient != nil {
		webSocketClient.Close()
	}
	return true
}

func (m *mattermostClient) findTeam() bool {
	var resp *model.Response
	if m.team, resp = m.client.GetTeamByName(gconfig.Mattermost.Team, ""); resp.Error != nil {
		log.Errorf("Failed to get team '%s', maybe we are not a member of this team.", gconfig.Mattermost.Team)
		return false
	}

	return true
}

func (m *mattermostClient) isMattermostUp() bool {
	if m.client == nil {
		return false
	}

	if _, resp := m.client.GetOldClientConfig(""); resp.Error != nil {
		log.Error("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		return false
	}

	if m.user == nil && !m.login() {
		return false
	}

	if m.team == nil && !m.findTeam() {
		return false
	}

	return true
}

func (m *mattermostClient) getChannelInfo(channelName string) *model.Channel {
	var channel *model.Channel
	var resp *model.Response

	if channel, resp = m.client.GetChannelByName(channelName, m.team.Id, ""); resp.Error != nil {
		log.Infof("Failed to get channel %s. Error: %s", channelName, resp.Error.Message)
		return nil
	}

	return channel
}

func (m *mattermostClient) getChannelInfoByID(channelID string) *model.Channel {
	var channel *model.Channel
	var resp *model.Response

	if channel, resp = m.client.GetChannel(channelID, ""); resp.Error != nil {
		log.Infof("Failed to get channel %s. Error: %s", channelID, resp.Error.Message)
		return nil
	}

	return channel
}

func (m *mattermostClient) createChannelIfNeeded(channelName string, channelDisplayName string,
	channelType string) bool {
	if m.client == nil || m.team == nil {
		log.Errorf("Client or team is nil, cannot create channel")
		return false
	}

	if chanInfo := m.getChannelInfo(channelName); chanInfo != nil {
		if chanInfo.DisplayName != channelDisplayName {
			if _, resp := mClient.client.UpdateChannel(&model.Channel{
				Name:        channelName,
				DisplayName: channelDisplayName,
				Purpose:     "IRC bridge with freenode %s channel",
				Type:        channelType,
				TeamId:      m.team.Id,
			}); resp.Error != nil {
				log.Errorf("Failed to update channel '%s': %s", channelName, resp.Error.Message)
				return false
			}
		}

		return true
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{
		Name:        channelName,
		DisplayName: channelDisplayName,
		Purpose:     "IRC bridge with freenode %s channel",
		Type:        channelType,
		TeamId:      m.team.Id,
	}

	if _, resp := m.client.CreateChannel(channel); resp.Error != nil {
		log.Errorf("Failed to create channel '%s': %s", channelName, resp.Error.Message)
		return false
	}

	return true
}

// handle websocket events from Mattermost
// return false when fatal error occur to close the channel
func (m *mattermostClient) handleWebSocketEvent(event *model.WebSocketEvent) bool {
	if event == nil {
		return false
	}

	log.Debugf("Event received type: %s, data %v", event.Event, event)
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return true
	}

	if _, ok := event.Data["post"]; !ok {
		log.Error("Malformed event found, 'post' key not found")
		return true
	}

	if _, ok := event.Data["sender_name"]; !ok {
		log.Error("Malformed event found, 'sender_name' key not found")
		return true
	}

	if _, ok := event.Data["channel_type"]; !ok {
		log.Error("Malformed event found, 'channel_type' key not found")
		return true
	}

	if event.Data["channel_type"] != "O" && event.Data["channel_type"] != "P" {
		log.Infof("Ignore event on non authorized channel_type %s", event.Data["channel_type"])
		return true
	}

	sender := event.Data["sender_name"].(string)

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {
		// ignore bot events and empty messages
		if post.UserId == mClient.user.Id || len(post.Message) == 0 {
			return true
		}

		// ignore webhook events (permits to break event loop on the user)
		if fromWebhook, ok := post.Props["from_webhook"]; ok && fromWebhook.(string) == "true" {
			return true
		}

		log.Debugf("Post received from other user: %v", post)

		if post.Message[0] != '!' {
			if err := m.sendIRCMessageToRabbit(post.Message, event.Broadcast.ChannelId, sender); err != nil {
				m.client.CreatePost(&model.Post{
					ChannelId: event.Broadcast.ChannelId,
					Message:   err.Error(),
				})
			}
			return true
		}

		// Ignore non command
		if len(post.Message) < 2 || post.Message[0] != '!' {
			return true
		}

		if err := m.sendCommandToRabbit(post.Message[1:], event.Broadcast.ChannelId,
			event.Broadcast.UserId); err != nil {
			m.client.CreatePost(&model.Post{
				ChannelId: event.Broadcast.ChannelId,
				Message:   err.Error(),
			})
		}
	}

	return true
}

func (m *mattermostClient) sendIRCMessageToRabbit(message string, channel string, sender string) error {
	if !gconfig.isAllowedIRCSender(sender) {
		return nil
	}

	chanInfos := m.getChannelInfoByID(channel)
	if chanInfos == nil {
		log.Errorf("Failed to find incoming user message channel %s, ignoring message.", channel)
		return nil
	}

	// Ignore non IRC channel forwarder
	if chanInfos.DisplayName[0:5] != "irc-#" {
		return nil
	}

	log.Debugf("Sender '%s' is allowed to send a message on IRC, forwarding message to IRC channel %s.",
		sender, chanInfos.DisplayName[4:])

	if !verifyPublisher() {
		log.Error("Failed to verify publisher, no command sent to broker. Notifying user.")
		return fmt.Errorf("%s: unable to send message to broker, message not sent", sender)
	}

	if !rabbitmqPublisher.Publish(
		&rabbitmq.IRCChatEvent{
			Message: message,
			Channel: chanInfos.DisplayName[4:],
			User:    sender,
		},
		"irc-chat",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			RoutingKey:    gconfig.Mattermost.IRCSenderRoutingKey,
			ExpirationMs:  300000,
		}) {
		log.Errorf("Failed to publish irc chat message to broker. Notifying user.")
		return fmt.Errorf("%s: unable to publish message to broker, message not sent", sender)
	}

	return nil
}

func (m *mattermostClient) sendCommandToRabbit(command string, channel string, user string) error {
	event := rabbitmq.CommandEvent{
		Command: command,
		Channel: channel,
		User:    user,
	}

	log.Infof("User %s sent command on channel %s: %s", event.User, event.Channel, event.Command)

	if !verifyPublisher() {
		log.Error("Failed to verify publisher, no command sent to broker. Notifying user.")
		return fmt.Errorf("%s: unable to send message to broker, command not sent", user)
	}

	if !verifyConsumer() {
		log.Error("Failed to verify consumer, no command sent to broker")
		return fmt.Errorf("%s: broker consumer has a problem, command not sent", user)
	}

	consumerCfg := gconfig.RabbitMQ.GetConsumer("commands")
	if consumerCfg == nil {
		log.Fatalf("RabbitMQ consumer configuration 'commands' not found, aborting.")
	}

	if !rabbitmqPublisher.Publish(
		&event,
		"command",
		&rabbitmq.EventOptions{
			CorrelationID: uuid.NewV4().String(),
			ReplyTo:       consumerCfg.RoutingKey,
			ExpirationMs:  300000,
		},
	) {
		log.Errorf("Failed to publish command to broker. Notifying user.")
		return fmt.Errorf("%s: unable to publish message to broker, command not sent", user)
	}

	return nil
}
