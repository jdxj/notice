package main

import (
	"context"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/robfig/cron/v3"

	"github.com/jdxj/notice/config"
)

type status struct {
	url     string
	publish time.Time
}

func NewRSS() *RSS {
	rss := &RSS{
		urls:  make(map[string]*status),
		cron:  cron.New(),
		parse: gofeed.NewParser(),
	}
	return rss
}

type RSS struct {
	urls  map[string]*status
	cron  *cron.Cron
	parse *gofeed.Parser
}

func (r *RSS) Start() {
	_, err := r.cron.AddFunc(config.RSS.Spec, func() {
		r.getURLs()
		r.run()
	})
	if err != nil {
		logger.Errorf("add func err: %s", err)
		return
	}
	r.cron.Start()
}

func (r *RSS) Stop() {
	<-r.cron.Stop().Done()
}

func (r *RSS) getURLs() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var urls []string
	err := db.WithContext(ctx).
		Table("rss_urls").
		Select("url").
		Find(&urls).Error
	if err != nil {
		logger.Errorf("get urls err: %s", err)
		return
	}

	for _, url := range urls {
		if _, ok := r.urls[url]; ok {
			continue
		}
		r.urls[url] = &status{url: url}
	}
}

func (r *RSS) run() {
	for url, status := range r.urls {
		feed, err := r.parse.ParseURL(url)
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

		if !status.publish.IsZero() && item.PublishedParsed.After(status.publish) {
			err := SendMessage(item.Link)
			if err != nil {
				logger.Errorf("send message err: %s", err)
			}
		}

		status.publish = *item.PublishedParsed
	}
}
