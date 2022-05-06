package github

import (
	"context"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"

	"github.com/jdxj/notice/config"
)

var (
	client *github.Client
)

func init() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken:  config.Github.PersonalAccessToken,
			TokenType:    "",
			RefreshToken: "",
			Expiry:       time.Time{},
		})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)
}

func GetProfile() {
	// releaseInfo, rsp, err := client.Repositories.
	// 	GetLatestRelease(context.Background(), "go-gost", "gost")
	// if err != nil {
	// 	log.Panicln(err)
	// }
	// fmt.Printf("%+v\n", rsp)
	// fmt.Printf("%+v\n", releaseInfo)

}

func Run() {

}

type issuesWatcher struct {
	UpdatedAt map[string]time.Time
}

func (iw *issuesWatcher) watchIssues() {
	interval := time.Second * time.Duration(config.Github.Interval)
	t := time.NewTimer(interval)
	defer t.Stop()

	for {
		<-t.C

		t.Reset(interval)
	}
}

func (iw *issuesWatcher) update() {

}
