package internal

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"gitlab.com/nerzhul/bot/rabbitmq"
	"golang.org/x/oauth2"
	"net/http"
)

type githubRepository struct {
	group string
	name  string
}

func checkGithubNewTags() bool {
	log.Infof("Checking new Github tags...")

	repositories, err := gDB.GetGithubConfiguredRepositories()
	if err != nil {
		log.Errorf("Failed to fetch Github configured repositories")
		return false
	}

	var httpClient *http.Client

	ctx := context.Background()

	// If there is a token defined, use it
	if len(gconfig.GithubToken) > 0 {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: gconfig.GithubToken},
		)
		httpClient = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(httpClient)

	for _, repo := range repositories {
		log.Infof("Fetching tags for github.com:%s/%s.git", repo.group, repo.name)

		tags, _, err := client.Repositories.ListTags(ctx, repo.group, repo.name, nil)
		if err != nil {
			log.Errorf("Failed to list Github repository tags: %s", err)
			return false
		}

		for _, t := range tags {
			registered, err := gDB.IsGithubRepositoryTagRegistered(repo.group, repo.name, *t.Name)
			if err != nil {
				return false
			}

			if !registered {
				// Unable to verify rabbitmq publisher, cancel this occurence
				if !verifyPublisher() {
					return false
				}

				publishAnnouncement(&rabbitmq.AnnouncementMessage{
					Message: *t.Name,
					What:    fmt.Sprintf("%s/%s", repo.group, repo.name),
					URL: fmt.Sprintf("https://github.com/%s/%s/releases/tag/%s",
						repo.group, repo.name, *t.Name),
				})

				if !gDB.RegisterRepositoryTag(repo.group, repo.name, *t.Name) {
					return false
				}

				// @TODO: send the notification
			}
		}
	}

	log.Infof("New Github tags fetch done.")
	return true
}
