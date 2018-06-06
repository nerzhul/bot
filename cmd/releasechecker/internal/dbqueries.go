package internal

const (
	getGithubRepositories    = `SELECT gh_group, gh_name FROM github_repositories`
	addGithubRepositoryQuery = `INSERT INTO github_repositories(gh_group, gh_name) VALUES ($1, $2) ON CONFLICT ` +
		`ON CONSTRAINT github_repositories_pkey DO NOTHING`
	addGithubRepositoryTag          = `INSERT INTO github_repository_tags (gh_group, gh_name, tag_name) VALUES ($1, $2, $3)`
	isGithubRepositoryTagRegistered = `SELECT EXISTS(SELECT 1 FROM github_repository_tags WHERE ` +
		`gh_group = $1 AND gh_name = $2 AND tag_name = $3);`
)
