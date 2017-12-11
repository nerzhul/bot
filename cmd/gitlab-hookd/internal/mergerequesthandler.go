package internal

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/satori/go.uuid"
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
	Changes          gitlabChanges                `json:"changes"`
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
		log.Error("Failed to read Gitlab Push event")
		return false
	}

	if !mrEvent.verifyEvent() {
		log.Error("Failed to verify Gitlab Tag Push event")
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
		rEvent := gitlabRabbitMQEvent{
			Message:     notificationMessage,
			Channel:     channel,
			User:        "",
			MessageType: "notice",
		}

		if !verifyPublisher() {
			log.Error("Failed to publish Gitlab Tag Push event")
			return false
		}

		rabbitmqPublisher.publish(&rEvent, uuid.NewV4().String())
	}
	return true
}
