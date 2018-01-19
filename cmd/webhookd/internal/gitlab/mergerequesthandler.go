package gitlab

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/common"
	"gitlab.com/nerzhul/bot/cmd/webhookd/internal/rabbitmq"
	"strings"
)

// swagger:parameters pushGitlabEvent
type gitlabMergeRequestEventParams struct {
	// in: body
	Body gitlabPushEvent
}

type gitlabMergeRequestEvent struct {
	ObjectKind       string                       `json:"object_kind"`
	User             gitlabUser                   `json:"user"`
	Project          gitlabProject                `json:"project"`
	Repository       gitlabRepository             `json:"repository"`
	ObjectAttributes gitlabMergeRequestAttributes `json:"object_attributes"`
	Labels           []gitlabLabel                `json:"labels"`
	// Changes field is not compliant with doc, ignore it
	//Changes          gitlabChanges                `json:"changes"`
}

func (gevent *gitlabMergeRequestEvent) verifyEvent() bool {
	if gevent.Project.PathWithNamespace == "" {
		return false
	}

	return true
}

func (gevent *gitlabMergeRequestEvent) toNotificationString() string {
	return fmt.Sprintf("[%s][MR !%d][%s] %s (%s)\n",
		gevent.ObjectAttributes.Target.PathWithNamespace,
		gevent.ObjectAttributes.IID,
		gevent.ObjectAttributes.Action,
		strings.Replace(gevent.ObjectAttributes.Title, "\n", "", -1),
		gevent.ObjectAttributes.URL)
}

func handleGitlabMergeRequest(c echo.Context) bool {
	mrEvent := gitlabMergeRequestEvent{}

	if !common.ReadJSONRequest(c.Request().Body, &mrEvent) {
		common.Log.Error("Failed to read Gitlab Merge Request event")
		return false
	}

	if !mrEvent.verifyEvent() {
		common.Log.Error("Failed to verify Gitlab Merge Request event")
		return false
	}

	channelsToPublish, exists := common.GConfig.GitlabProjectsMapping[mrEvent.Project.PathWithNamespace]
	if !exists {
		common.Log.Warningf("Received hook from project %s but not channel mapped.",
			mrEvent.Project.PathWithNamespace)
		return true
	}

	notificationMessage := mrEvent.toNotificationString()

	for _, channel := range channelsToPublish {
		rEvent := bot.CommandResponse{
			Message:     notificationMessage,
			Channel:     channel,
			User:        "",
			MessageType: "notice",
		}

		if !rabbitmq.VerifyPublisher() {
			common.Log.Error("Failed to publish Gitlab Merge Request event")
			return false
		}

		rabbitmq.Publisher.Publish(
			&rEvent,
			"gitlab-event",
			&bot.EventOptions{
				CorrelationID: uuid.NewV4().String(),
				ExpirationMs:  300000,
			},
		)
	}
	return true
}
