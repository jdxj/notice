package github

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
	"github.com/jdxj/notice/subscription"
	"github.com/jdxj/notice/subscription/telegram"
	"github.com/jdxj/notice/util"
	"github.com/jdxj/notice/util/db"
)

var (
	repoIDMap = make(map[string]time.Time)
	crontab   = cron.New()
	client    = newGithubClient()
)

func newGithubClient() *github.Client {
	ts := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: config.Github.PersonalAccessToken,
	})
	tc := oauth2.NewClient(context.Background(), ts)
	return github.NewClient(tc)
}

func init() {
	restore()
	start()
}

func repoID(owner, repo string) string {
	return fmt.Sprintf("%s/%s", owner, repo)
}

func restore() {
	ctx, cancel := util.WithTimeout()
	defer cancel()

	var rows []entity
	err := db.WithContext(ctx).
		Model(entity{}).
		Find(&rows).Error
	if err != nil {
		logger.Panicf("restore github repos: %s", err)
	}

	for _, row := range rows {
		repoIDMap[repoID(row.Owner, row.Repo)] = time.Time{}
	}
}

func start() {
	_, err := crontab.AddFunc(config.Github.Spec, checkPublish)
	if err != nil {
		logger.Panicf("add func err: %s", err)
	}

	crontab.Start()
	logger.Infof("github started")
}

func checkPublish() {
	for repoID, lastPublish := range repoIDMap {
		ur := strings.Split(repoID, "/")
		releases, _, err := client.Repositories.ListReleases(context.Background(), ur[0], ur[1], nil)
		if err != nil {
			logger.Errorf("list releases err: %s", err)
			continue
		}

		if len(releases) == 0 {
			logger.Warnf("releases not found: %s", repoID)
			continue
		}

		release := releases[0]
		if !lastPublish.IsZero() && release.GetPublishedAt().After(lastPublish) {
			err := telegram.SendMessage(fmt.Sprintf("%s is updated: %s", repoID, release.GetTagName()))
			if err != nil {
				logger.Errorf("send message err: %s", err)
			}
		}

		repoIDMap[repoID] = release.GetPublishedAt().Time
	}
}

type AddReq struct {
	Owner string
	Repo  string
}

func Add(ctx context.Context, req *AddReq) error {
	key := repoID(req.Owner, req.Repo)
	if !repoIDMap[key].IsZero() {
		return fmt.Errorf("%w: %s", subscription.ErrKeyAlreadyExisted, key)
	}
	repoIDMap[key] = time.Time{}

	return db.WithContext(ctx).
		Create(entity{Owner: req.Owner, Repo: req.Repo}).
		Error
}
