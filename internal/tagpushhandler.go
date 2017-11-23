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
	UserId           uint64           `json:"user_id"`
	UserName         string           `json:"user_name"`
	UserUsername     string           `json:"user_username"`
	UserAvatar       string           `json:"user_avatar"`
	ProjectId        uint64           `json:"project_id"`
	Commits          []gitlabCommit   `json:"commits"`
	TotalCommitCount uint64           `json:"total_commits_count"`
}

func handleGitlabTagPush(c echo.Context) bool {
	tagPushEvent := gitlabTagPushEvent{}

	if !ReadJsonRequest(c.Request().Body, &tagPushEvent) {
		return false
	}

	if tagPushEvent.TotalCommitCount == 0 || tagPushEvent.Project.PathWithNamespace == "" {
		return false
	}

	channelsToPublish, kExist := gconfig.ProjectsMapping[tagPushEvent.Project.PathWithNamespace]
	if !kExist {
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
			return false
		}

		rabbitmqPublisher.Publish(&rEvent, uuid.NewV4().String())
	}
	return true
}
