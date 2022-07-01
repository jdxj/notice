package github

type entity struct {
	ID    int
	Owner string
	Repo  string
}

func (entity) TableName() string {
	return "github_repos"
}
