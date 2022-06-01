package github

import (
	"context"
	"log"
	"time"

	"github.com/google/go-github/v44/github"
	"golang.org/x/oauth2"

	"github.com/jdxj/notice/config"
)

func New() *Watcher {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken:  config.Github.PersonalAccessToken,
			TokenType:    "",
			RefreshToken: "",
			Expiry:       time.Time{},
		})
	tc := oauth2.NewClient(ctx, ts)
	return &Watcher{
		client: github.NewClient(tc),
	}
}

type Watcher struct {
	client *github.Client
}

func (w *Watcher) Name() string {
	return "github"
}

func (w *Watcher) Watch(ctx context.Context) (string, bool, error) {
	tags, _, err := w.client.Repositories.ListTags(ctx, "go-gost", "gost", nil)
	if err != nil {
		return "", false, err
	}
	for _, v := range tags {
		log.Printf("%+v", v.GetName())
	}
	return "", false, nil
}
