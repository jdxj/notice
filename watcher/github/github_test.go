package github

import (
	"context"
	"testing"
)

func TestWatcher_Watch(t *testing.T) {
	w := New()
	w.Watch(context.Background())
}
