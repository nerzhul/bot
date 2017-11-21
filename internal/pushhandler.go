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
	ObjectKind   string `json:"object_kind"`
	Before       string `json:"before"`
	After        string `json:"after"`
	Ref          string `json:"ref"`
	CheckoutSha  string `json:"checkout_sha"`
	UserId       uint64 `json:"user_id"`
	UserName     string `json:"user_name"`
	UserUsername string `json:"user_username"`
	UserAvatar   string `json:"user_avatar"`
	ProjectId    uint64 `json:"project_id"`
	Project      struct {
		Id                uint64 `json:"id"`
		Name              string `json:"name"`
		Description       string `json:"description"`
		WebUrl            string `json:"web_url"`
		AvatarUrl         string `json:"avatar_url"`
		GitSSHUrl         string `json:"git_ssh_url"`
		GitHTTPUrl        string `json:"git_http_url"`
		Namespace         string `json:"namespace"`
		VisibilityLevel   uint16 `json:"visibility_level"`
		PathWithNamespace string `json:"path_with_namespace"`
		DefaultBranch     string `json:"default_branch"`
		Homepage          string `json:"homepage"`
		Url               string `json:"url"`
		SSHUrl            string `json:"ssh_url"`
		HTTPUrl           string `json:"http_url"`
	} `json:"project"`
	Repository struct {
		Name            string `json:"name"`
		Url             string `json:"url"`
		Description     string `json:"description"`
		Homepage        string `json:"homepage"`
		GitSSHUrl       string `json:"git_ssh_url"`
		GitHTTPUrl      string `json:"git_http_url"`
		VisibilityLevel uint16 `json:"visibility_level"`
	} `json:"repository"`
	Commits []struct {
		Id        string `json:"id"`
		Message   string `json:"message"`
		Timestamp string `json:"timestamp"` // @TODO: convert this
		Url       string `json:"url"`
		Author    struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		} `json:"author"`
		Added    []string `json:"added"`
		Modified []string `json:"modified"`
		Removed  []string `json:"removed"`
	} `json:"commits"`
	TotalCommitCount uint64 `json:"total_commits_count"`
}

func handleGitlabPush(c echo.Context) bool {
	pushEvent := gitlabPushEvent{}

	if !ReadJsonRequest(c.Request().Body, &pushEvent) {
		return false
	}

	if pushEvent.TotalCommitCount == 0 || pushEvent.Project.PathWithNamespace == "" {
		return false
	}

	channelsToPublish, kExist := gconfig.ProjectsMapping[pushEvent.Project.PathWithNamespace]
	if !kExist {
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

	notificationMessage += strings.Replace(lastCommit.Message, "\n", "", -1) + " (" + lastCommit.Url + ")\n"

	rEvent := gitlabRabbitMQEvent{Message: notificationMessage, Channels: channelsToPublish}

	if !verifyPublisher() {
		return false
	}

	rabbitmqPublisher.Publish(&rEvent, uuid.NewV4().String())
	return true
}
