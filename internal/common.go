package internal

type gitlabRepository struct {
	Name            string `json:"name"`
	Url             string `json:"url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	GitSSHUrl       string `json:"git_ssh_url"`
	GitHTTPUrl      string `json:"git_http_url"`
	VisibilityLevel uint16 `json:"visibility_level"`
}

type gitlabProject struct {
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
}

type gitlabCommit struct {
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
}
