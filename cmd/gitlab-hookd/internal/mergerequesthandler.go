package internal

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
	"gitlab.com/nerzhul/bot"
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

	if !readJSONRequest(c.Request().Body, &mrEvent) {
		log.Error("Failed to read Gitlab Merge Request event")
		return false
	}

	if !mrEvent.verifyEvent() {
		log.Error("Failed to verify Gitlab Merge Request event")
		return false
	}

	channelsToPublish, exists := gconfig.ProjectsMapping[mrEvent.Project.PathWithNamespace]
	if !exists {
		log.Warningf("Received hook from project %s but not channel mapped.",
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

		if !verifyPublisher() {
			log.Error("Failed to publish Gitlab Merge Request event")
			return false
		}

		rabbitmqPublisher.Publish(
			&rEvent,
			"gitlab-event",
			uuid.NewV4().String(),
			"",
			300000,
		)
	}
	return true
}
