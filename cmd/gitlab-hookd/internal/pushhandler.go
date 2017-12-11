package internal

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"strconv"
	"strings"
)

// swagger:parameters pushGitlabEvent
type gitlabPushEventParams struct {
	// in: body
	Body gitlabPushEvent
}

type gitlabPushEvent struct {
	ObjectKind       string           `json:"object_kind"`
	Before           string           `json:"before"`
	After            string           `json:"after"`
	Ref              string           `json:"ref"`
	CheckoutSha      string           `json:"checkout_sha"`
	UserID           uint64           `json:"user_id"`
	UserName         string           `json:"user_name"`
	UserUsername     string           `json:"user_username"`
	UserAvatar       string           `json:"user_avatar"`
	ProjectID        uint64           `json:"project_id"`
	Project          gitlabProject    `json:"project"`
	Repository       gitlabRepository `json:"repository"`
	Commits          []gitlabCommit   `json:"commits"`
	TotalCommitCount uint64           `json:"total_commits_count"`
}

func handleGitlabPush(c echo.Context) bool {
	pushEvent := gitlabPushEvent{}

	if !readJSONRequest(c.Request().Body, &pushEvent) {
		return false
	}

	if pushEvent.TotalCommitCount == 0 || pushEvent.Project.PathWithNamespace == "" {
		return false
	}

	channelsToPublish, exists := gconfig.ProjectsMapping[pushEvent.Project.PathWithNamespace]
	if !exists {
		log.Warningf("Received hook from project %s but not channel mapped.",
			pushEvent.Project.PathWithNamespace)
		return true
	}

	pushEvent.Ref = strings.Replace(pushEvent.Ref, "refs/heads/", "", -1)

	var notificationMessage string
	notificationMessage += "[" + pushEvent.Project.PathWithNamespace + "][" + pushEvent.Ref + "] " +
		pushEvent.UserName + " pushed " + strconv.FormatUint(pushEvent.TotalCommitCount, 10) + " commit"
	if pushEvent.TotalCommitCount > 1 {
		notificationMessage += "s"
	}

	notificationMessage += ". "

	if pushEvent.TotalCommitCount > 1 {
		notificationMessage += "Last: "
	}

	lastCommit := &pushEvent.Commits[0]

	notificationMessage += strings.Replace(lastCommit.Message, "\n", "", -1) + " (" + lastCommit.URL + ")\n"

	for _, channel := range channelsToPublish {
		rEvent := gitlabRabbitMQEvent{
			Message:     notificationMessage,
			Channel:     channel,
			User:        "",
			MessageType: "notice",
		}

		if !verifyPublisher() {
			return false
		}

		rabbitmqPublisher.publish(&rEvent, uuid.NewV4().String())
	}
	return true
}
