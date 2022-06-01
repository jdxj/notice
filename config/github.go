package config

type github struct {
	// https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
	PersonalAccessToken string `mapstructure:"personal_access_token"`

	// 要 watch 的事件
	// 每个仓库使用 owner/repo 格式
	Issues         []string `mapstructure:"issues"`
	PullRequests   []string `mapstructure:"pull_requests"`
	Releases       []string `mapstructure:"releases"`
	Discussions    []string `mapstructure:"discussions"`
	SecurityAlerts []string `mapstructure:"security_alerts"`
}
