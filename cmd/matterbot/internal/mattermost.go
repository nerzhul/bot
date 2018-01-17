package internal

import (
	"github.com/mattermost/mattermost-server/model"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot"
	"strings"
	"time"
)

type mattermostClient struct {
	client *model.Client4
	user   *model.User
	team   *model.Team
}

var mClient mattermostClient

func runMattermostClient() {
	for {
		mClient.init()
		if mClient.isMattermostUp() {
			mClient.createChannelIfNeeded(gconfig.Mattermost.TwitterChannel, model.CHANNEL_PRIVATE)
		}

		mClient.run()

		mClient.deinit()
		log.Warning("Connection to mattermost lost, retrying in 60s...")
		time.Sleep(time.Second * 60)
	}
}

func (m *mattermostClient) init() {
	m.client = model.NewAPIv4Client(gconfig.Mattermost.URL)
}

func (m *mattermostClient) deinit() {
	mClient.client = nil
	mClient.user = nil
	mClient.team = nil
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
			m.handleWebSocketResponse(resp)
		}
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

func (m *mattermostClient) createChannelIfNeeded(channelName string, channelType string) bool {
	if m.client == nil || m.team == nil {
		log.Errorf("Client or team is nil, cannot create channel")
		return false
	}

	if _, resp := m.client.GetChannelByName(channelName, m.team.Id, ""); resp.Error != nil {
		log.Infof("Failed to get channels %s. Error: %s. Trying to create channel",
			channelName, resp.Error.Message)
	} else {
		// Channel already exists, ignore
		return true
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = channelName
	channel.DisplayName = channelName
	channel.Purpose = ""
	channel.Type = channelType
	channel.TeamId = m.team.Id
	if _, resp := m.client.CreateChannel(channel); resp.Error != nil {
		log.Errorf("Failed to create channel '%s'", channelName)
		return false
	}

	return true
}

func (m *mattermostClient) handleWebSocketResponse(event *model.WebSocketEvent) {
	if event == nil {
		return
	}

	log.Debugf("Event received type: %s", event.Event)
	if event.Event != model.WEBSOCKET_EVENT_POSTED {
		return
	}

	post := model.PostFromJson(strings.NewReader(event.Data["post"].(string)))
	if post != nil {
		log.Debugf("Post received: %s", post)
		// ignore bot events
		if post.UserId == mClient.user.Id {
			return
		}

		// Ignore non command
		if len(post.Message) < 2 || post.Message[0] != '!' {
			return
		}

		event := bot.CommandEvent{
			Command: post.Message[1:],
			Channel: event.Broadcast.ChannelId,
			User:    event.Broadcast.UserId,
		}

		log.Infof("User %s sent command on channel %s: %s", event.User, event.Channel, event.Command)

		if !verifyPublisher() {
			log.Error("Failed to verify publisher, no command sent to broker")
			return
		}

		if !verifyConsumer() {
			log.Error("Failed to verify consumer, no command sent to broker")
			return
		}

		consumerCfg := gconfig.RabbitMQ.GetConsumer("commands")
		if consumerCfg == nil {
			log.Fatalf("RabbitMQ consumer configuration 'commands' not found, aborting.")
		}

		rabbitmqPublisher.Publish(
			&event,
			"command",
			&bot.EventOptions{
				CorrelationID: uuid.NewV4().String(),
				ReplyTo:       consumerCfg.RoutingKey,
				ExpirationMs:  300000,
			},
		)
	}
}
