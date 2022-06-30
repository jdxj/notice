package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/go-github/v45/github"
	"github.com/robfig/cron/v3"
	"golang.org/x/oauth2"

	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/logger"
)

func NewGithub() *Github {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: config.Github.PersonalAccessToken,
	})
	tc := oauth2.NewClient(context.Background(), ts)

	g := &Github{
		repos:  make(map[string]*status),
		cron:   cron.New(),
		client: github.NewClient(tc),
	}
	return g
}

type Github struct {
	repos  map[string]*status
	cron   *cron.Cron
	client *github.Client
}

func (g *Github) Start() {
	_, err := g.cron.AddFunc(config.Github.Spec, func() {
		g.getRepos()
		g.run()
	})
	if err != nil {
		logger.Errorf("add func err: %s", err)
		return
	}
	g.cron.Start()
}

func (g *Github) Stop() {
	<-g.cron.Stop().Done()
}

type uniqueRepo struct {
	Owner string
	Repo  string
}

func (g *Github) getRepos() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var repos []uniqueRepo
	err := db.WithContext(ctx).
		Table("github").
		Select("owner, repo").
		Find(&repos).Error
	if err != nil {
		logger.Errorf("get repos err: %s", err)
		return
	}

	for _, v := range repos {
		key := fmt.Sprintf("%s/%s", v.Owner, v.Repo)
		if _, ok := g.repos[key]; ok {
			continue
		}
		g.repos[key] = &status{key: key}
	}
}

func (g *Github) run() {
	for key, status := range g.repos {
		ur := strings.Split(key, "/")
		releases, _, err := g.client.Repositories.ListReleases(context.Background(), ur[0], ur[1], nil)
		if err != nil {
			logger.Errorf("list releases err: %s", err)
			continue
		}

		if len(releases) == 0 {
			logger.Warnf("releases not found: %s", key)
			continue
		}

		release := releases[0]
		if !status.publish.IsZero() && release.GetPublishedAt().After(status.publish) {
			err := SendMessage(fmt.Sprintf("%s is updated: %s", key, release.GetTagName()))
			if err != nil {
				logger.Errorf("send message err: %s", err)
			}
		}

		status.publish = release.GetPublishedAt().Time
	}
}
