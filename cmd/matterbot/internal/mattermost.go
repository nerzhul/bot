package internal

import (
	"github.com/mattermost/mattermost-server/model"
)

type mattermostClient struct {
	client *model.Client4
	user   *model.User
	team   *model.Team
}

var mClient mattermostClient

func runMattermostClient() {
	mClient.init()
}

func (m *mattermostClient) init() {
	m.client = model.NewAPIv4Client(gconfig.Mattermost.URL)
}

func (m *mattermostClient) login() bool {
	if user, resp := m.client.Login(gconfig.Mattermost.Email, gconfig.Mattermost.Password); resp.Error != nil {
		log.Error("There was a problem logging into the Mattermost server. Ensure login is correct.")
		return false
	} else {
		m.user = user
	}

	// Lets start listening to some channels via the websocket!
	webSocketClient, err := model.NewWebSocketClient4(gconfig.Mattermost.WsURL, m.client.AuthToken)
	if err != nil {
		log.Error("Failed to connect to Mattermost websocket.")
		return false
	}

	webSocketClient.Listen()

	go func() {
		for {
			select {
			case resp := <-webSocketClient.EventChannel:
				m.handleWebSocketResponse(resp)
			}
		}
	}()

	// You can block forever with
	select {}

	return true
}

func (m *mattermostClient) findTeam() bool {
	if team, resp := m.client.GetTeamByName(gconfig.Mattermost.Team, ""); resp.Error != nil {
		log.Errorf("Failed to get team '%s', maybe we are not a member of this team.", gconfig.Mattermost.Team)
		return false
	} else {
		m.team = team
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

func (m *mattermostClient) handleWebSocketResponse(event *model.WebSocketEvent) {
	// @TODO
}
