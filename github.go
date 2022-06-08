package main

import (
	"context"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"

	"github.com/jdxj/notice/config"
)

func NewGithub() *Github {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: config.Github.PersonalAccessToken,
	})
	tc := oauth2.NewClient(context.Background(), ts)

	g := &Github{
		client: github.NewClient(tc),
	}
	return g
}

type Github struct {
	client *github.Client
}

func (g *Github) T() {
	g.client.Repositories.ListTags()
}
