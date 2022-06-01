package scheduler

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron/v3"

	"github.com/jdxj/notice/telegram"
	"github.com/jdxj/notice/watcher"
)

func New() *Scheduler {
	return &Scheduler{
		cron:     cron.New(),
		messages: make(chan string, 10),
		watchers: make(map[string]watcher.Watcher),
	}
}

type Scheduler struct {
	cron     *cron.Cron
	messages chan string
	watchers map[string]watcher.Watcher
}

func (s *Scheduler) Register(w watcher.Watcher) {
	if w == nil {
		return
	}
	if _, ok := s.watchers[w.Name()]; ok {
		return
	}
	_, _ = s.cron.AddFunc("*/5 * * * *", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()

		msg, changed, err := w.Watch(ctx)
		if err != nil {
			log.Printf("%s watch err: %s", w.Name(), err)
			return
		}

		if changed {
			s.messages <- msg
		}
	})
}

func (s *Scheduler) Run() {
	go s.sendMessages()
	s.cron.Run()
}

func (s *Scheduler) sendMessages() {
	for msg := range s.messages {
		err := telegram.SendMessage(msg)
		if err != nil {
			log.Printf("sendMessage err: %s", err)
		}
	}
}
