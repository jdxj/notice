package rss

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"

	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/logger"
	"github.com/jdxj/notice/subscription/telegram"
	"github.com/jdxj/notice/util"
	"github.com/jdxj/notice/util/db"
)

var (
	urlMap  = make(map[string]time.Time)
	crontab = cron.New()
	parser  = gofeed.NewParser()
)

func init() {
	restore()
	start()
}

func restore() {
	ctx, cancel := util.WithTimeout()
	defer cancel()

	var rows []entity
	err := db.WithContext(ctx).
		Model(entity{}).
		Find(&rows).Error
	if err != nil {
		logger.Panicf("restore rss urls: %s", err)
	}

	for _, row := range rows {
		urlMap[row.URL] = time.Time{}
	}
}

func start() {
	_, err := crontab.AddFunc(config.RSS.Spec, checkPublish)
	if err != nil {
		logger.Panicf("add func err: %s", err)
	}

	crontab.Start()
	logger.Infof("rss started")
}

func checkPublish() {
	for url, lastPublish := range urlMap {
		feed, err := parser.ParseURL(url)
		if err != nil {
			logger.Errorf("parse url err: %s", err)
			continue
		}

		if len(feed.Items) == 0 {
			logger.Warnf("update not found: %s", url)
			continue
		}

		item := feed.Items[0]
		if item.PublishedParsed == nil {
			now := time.Now()
			item.PublishedParsed = &now
			logger.Warnf("published time not found: %s", url)
		}

		if !lastPublish.IsZero() && item.PublishedParsed.After(lastPublish) {
			err := telegram.SendMessage(item.Link)
			if err != nil {
				logger.Errorf("send message err: %s", err)
			}
		}

		urlMap[url] = *item.PublishedParsed
	}
}

// todo: validate
type AddReq struct {
	URL string
}

func Add(ctx context.Context, req *AddReq) error {
	urlMap[req.URL] = time.Time{}

	return db.WithContext(ctx).
		Create(entity{URL: req.URL}).
		Error
}
