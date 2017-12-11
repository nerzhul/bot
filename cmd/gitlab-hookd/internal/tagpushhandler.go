package internal

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"strings"
)

// swagger:parameters tagPushGitlabEvent
type gitlabTagPushEventParams struct {
	// in: body
	Body gitlabTagPushEvent
}

type gitlabTagPushEvent struct {
	Project          gitlabProject    `json:"project"`
	Repository       gitlabRepository `json:"repository"`
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
	Commits          []gitlabCommit   `json:"commits"`
	TotalCommitCount uint64           `json:"total_commits_count"`
}

func (gevent *gitlabTagPushEvent) verifyEvent() bool {
	if gevent.TotalCommitCount != 0 || gevent.Project.PathWithNamespace == "" {
		return false
	}

	return true
}

func (gevent *gitlabTagPushEvent) toNotificationString() string {
	gevent.Ref = strings.Replace(gevent.Ref, "refs/tags/", "", -1)

	return "[" + gevent.Project.PathWithNamespace + "] " + gevent.UserName +
		" pushed tag " + gevent.Ref + ".\n"
}

func handleGitlabTagPush(c echo.Context) bool {
	tagPushEvent := gitlabTagPushEvent{}

	if !readJSONRequest(c.Request().Body, &tagPushEvent) {
		log.Error("Failed to read Gitlab Tag Push event")
		return false
	}

	if !tagPushEvent.verifyEvent() {
		log.Error("Failed to verify Gitlab Tag Push event")
		return false
	}

	channelsToPublish, exists := gconfig.ProjectsMapping[tagPushEvent.Project.PathWithNamespace]
	if !exists {
		log.Warningf("Received hook from project %s but not channel mapped.",
			tagPushEvent.Project.PathWithNamespace)
		return true
	}

	tagPushEvent.Ref = strings.Replace(tagPushEvent.Ref, "refs/heads/", "", -1)

	var notificationMessage string
	notificationMessage += "[" + tagPushEvent.Project.PathWithNamespace + "] " + tagPushEvent.UserName +
		" pushed tag " + tagPushEvent.Ref + ".\n"

	for _, channel := range channelsToPublish {
		rEvent := gitlabRabbitMQEvent{
			Message:     notificationMessage,
			Channel:     channel,
			User:        "",
			MessageType: "notice",
		}

		if !verifyPublisher() {
			log.Error("Failed to publish Gitlab Tag Push event notification")
			return false
		}

		rabbitmqPublisher.publish(&rEvent, uuid.NewV4().String())
	}
	return true
}
