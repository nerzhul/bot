package internal

import (
	"fmt"
	"github.com/heroku/docker-registry-client/registry"
	"gitlab.com/nerzhul/bot/rabbitmq"
)

type dockerHubImage struct {
	group string
	name  string
}

func checkDockerHubNewTags() bool {
	log.Infof("Checking new DockerHub tags...")

	images, err := gDB.GetDockerHubConfiguredImages()
	if err != nil {
		log.Errorf("Failed to fetch DockerHub configured images")
		return false
	}

	huburl := "https://registry-1.docker.io/"
	username := "" // anonymous
	password := "" // anonymous
	hub, err := registry.New(huburl, username, password)
	// Change the registry logger
	hub.Logf = log.Infof

	if err != nil {
		log.Errorf("Failed to create DockerHub registry client: %s", err)
		return false
	}

	for _, image := range images {
		log.Infof("Fetching DockerHub image %s/%s tags", image.group, image.name)
		tags, err := hub.Tags(fmt.Sprintf("%s%%2F%s", image.group, image.name))
		if err != nil {
			log.Errorf("Failed to list DockerHub image tags: %s", err)
			return false
		}

		for _, t := range tags {
			registered, err := gDB.IsDockerHubImageTagRegistered(image.group, image.name, t)
			if err != nil {
				return false
			}

			if !registered {
				// Unable to verify rabbitmq publisher, cancel this occurence
				if !verifyPublisher() {
					return false
				}

				publishAnnouncement(&rabbitmq.AnnouncementMessage{
					Message: t,
					What:    fmt.Sprintf("%s/%s", image.group, image.name),
					URL: fmt.Sprintf("https://hub.docker.com/r/%s/%s/tags",
						image.group, image.name),
				})

				if !gDB.RegisterDockerHubImageTag(image.group, image.name, t) {
					return false
				}
			}
		}

	}

	log.Infof("New DockerHub tags fetch done.")
	return true
}
