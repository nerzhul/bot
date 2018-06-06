package internal

type githubRepository struct {
	group string
	name  string
}

func verifyGithubNewTags() bool {
	_, err := gDB.GetGithubConfiguredRepositories()
	if err != nil {
		log.Errorf("Failed to fetch Github configured repositories")
		return false
	}

	return true
}
