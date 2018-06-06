package internal

const (
	getGithubRepositories    = `SELECT gh_group, gh_name FROM github_repositories`
	addGithubRepositoryQuery = `INSERT INTO github_repositories(gh_group, gh_name) VALUES ($1, $2) ON CONFLICT ` +
		`ON CONSTRAINT github_repositories_pkey DO NOTHING`
)
