package watcher

import "context"

type Watcher interface {
	Name() string
	Watch(context.Context) (string, bool, error)
}
