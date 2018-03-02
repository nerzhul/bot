package gitlab

import (
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	myrabbitmq "gitlab.com/nerzhul/bot/cmd/webhookd/internal/rabbitmq"
	"gitlab.com/nerzhul/bot/rabbitmq"
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
	if gevent.Project.PathWithNamespace == "" {
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

	if !common.ReadJSONRequest(c.Request().Body, &tagPushEvent) {
		common.Log.Error("Failed to read Gitlab Tag Push event")
		return false
	}

	if !tagPushEvent.verifyEvent() {
		common.Log.Error("Failed to verify Gitlab Tag Push event")
		return false
	}

	channelsToPublish, exists := common.GConfig.GitlabProjectsMapping[tagPushEvent.Project.PathWithNamespace]
	if !exists {
		common.Log.Warningf("Received hook from project %s but not channel mapped.",
			tagPushEvent.Project.PathWithNamespace)
		return true
	}

	notificationMessage := tagPushEvent.toNotificationString()

	for _, channel := range channelsToPublish {
		rEvent := rabbitmq.CommandResponse{
			Message:     notificationMessage,
			Channel:     channel,
			User:        "",
			MessageType: "notice",
		}

		if !myrabbitmq.VerifyPublisher() {
			common.Log.Error("Failed to publish Gitlab Tag Push event notification")
			return false
		}

		myrabbitmq.Publisher.Publish(
			&rEvent,
			"gitlab-event",
			&rabbitmq.EventOptions{
				CorrelationID: uuid.NewV4().String(),
				ExpirationMs:  300000,
			},
		)
	}
	return true
}
