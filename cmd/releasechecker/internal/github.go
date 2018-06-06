package internal

import (
	"context"
	"github.com/google/go-github/github"
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

	client := github.NewClient(nil)

	for _, repo := range repositories {
		log.Infof("Fetching tags for github.com:%s/%s.git", repo.group, repo.name)

		tags, _, err := client.Repositories.ListTags(context.Background(), repo.group, repo.name, nil)
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
