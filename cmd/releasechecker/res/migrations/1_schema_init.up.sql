CREATE TABLE github_repositories (
  gh_group TEXT,
  gh_name TEXT,
  PRIMARY KEY (gh_group, gh_name)
);

CREATE TABLE github_repository_tags (
  gh_group TEXT,
  gh_name TEXT,
  tag_name TEXT,
  FOREIGN KEY (gh_group, gh_name) REFERENCES github_repositories(gh_group, gh_name),
  PRIMARY KEY (gh_group, gh_name, tag_name)
);