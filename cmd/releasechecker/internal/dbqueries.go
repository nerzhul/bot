package internal

const (
	// Github
	getGithubRepositories    = `SELECT gh_group, gh_name FROM github_repositories`
	addGithubRepositoryQuery = `INSERT INTO github_repositories(gh_group, gh_name) VALUES ($1, $2) ON CONFLICT ` +
		`ON CONSTRAINT github_repositories_pkey DO NOTHING`
	addGithubRepositoryTag          = `INSERT INTO github_repository_tags (gh_group, gh_name, tag_name) VALUES ($1, $2, $3)`
	isGithubRepositoryTagRegistered = `SELECT EXISTS(SELECT 1 FROM github_repository_tags WHERE ` +
		`gh_group = $1 AND gh_name = $2 AND tag_name = $3);`

	// DockerHub
	getDockerHubImages     = `SELECT dh_group, dh_name FROM dockerhub_images`
	addDockerHubImageQuery = `INSERT INTO dockerhub_images(dh_group, dh_name) VALUES ($1, $2) ON CONFLICT ` +
		`ON CONSTRAINT dockerhub_images_pkey DO NOTHING`
	addDockerHubImageTag          = `INSERT INTO dockerhub_image_tags (dh_group, dh_name, tag_name) VALUES ($1, $2, $3)`
	isDockerHubImageTagRegistered = `SELECT EXISTS(SELECT 1 FROM dockerhub_image_tags WHERE ` +
		`dh_group = $1 AND dh_name = $2 AND tag_name = $3);`
)
