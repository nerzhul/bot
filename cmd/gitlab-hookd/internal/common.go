package internal

type gitlabUser struct {
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
}

type gitlabRepository struct {
	Name            string `json:"name"`
	URL             string `json:"url"`
	Description     string `json:"description"`
	Homepage        string `json:"homepage"`
	GitSSHUrl       string `json:"git_ssh_url"`
	GitHTTPUrl      string `json:"git_http_url"`
	VisibilityLevel uint16 `json:"visibility_level"`
}

type gitlabProject struct {
	ID                uint64 `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	WebURL            string `json:"web_url"`
	AvatarURL         string `json:"avatar_url"`
	GitSSHURL         string `json:"git_ssh_url"`
	GitHTTPURL        string `json:"git_http_url"`
	Namespace         string `json:"namespace"`
	VisibilityLevel   uint16 `json:"visibility_level"`
	PathWithNamespace string `json:"path_with_namespace"`
	DefaultBranch     string `json:"default_branch"`
	Homepage          string `json:"homepage"`
	URL               string `json:"url"`
	SSHUrl            string `json:"ssh_url"`
	HTTPUrl           string `json:"http_url"`
}

type gitlabCommit struct {
	ID        string `json:"id"`
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"` // @TODO: convert this
	URL       string `json:"url"`
	Author    struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	} `json:"author"`
	Added    []string `json:"added"`
	Modified []string `json:"modified"`
	Removed  []string `json:"removed"`
}

type gitlabLabel struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Color       string `json:"string"`
	ProjectID   int    `json:"project_id"`
	CreatedAt   string `json:"created_at"` // @TODO convert this
	UpdatedAt   string `json:"updated_at"` // @TODO convert this
	Template    bool   `json:"template"`
	Description string `json:"description"`
	Type        string `json:"type"`
	GroupID     int    `json:"group_id"`
}

type gitlabChanges struct {
	UpdatedByID []int    `json:"updated_by_id"`
	UpdatedAt   []string `json:"updated_at"` // @TODO convert this
	Labels      struct {
		Previous []gitlabLabel `json:"previous"`
		Current  []gitlabLabel `json:"current"`
	} `json:"labels"`
}

type gitlabMergeRequestAttributes struct {
	ID              int           `json:"id"`
	TargetBranch    string        `json:"target_branch"`
	SourceBranch    string        `json:"source_branch"`
	SourceProjectID int           `json:"source_project_id"`
	AuthorID        int           `json:"author_id"`
	AssigneeID      int           `json:"assignee_id"`
	Title           string        `json:"title"`
	CreatedAt       string        `json:"created_at"` // @TODO convert this
	UpdatedAt       string        `json:"updated_at"` // @TODO convert this
	MilestoneID     int           `json:"milestone_id"`
	State           string        `json:"state"`
	MergeStatus     string        `json:"merge_status"`
	TargetProjectID int           `json:"target_project_id"`
	IID             int           `json:"iid"`
	Description     string        `json:"description"`
	Source          gitlabProject `json:"source"`
	Target          gitlabProject `json:"target"`
	LastCommit      gitlabCommit  `json:"last_commit"`
	WorkInProgress  bool          `json:"work_in_progress"`
	URL             string        `json:"url"`
	Action          string        `json:"action"`
	Assignee        gitlabUser    `json:"assignee"`
}
