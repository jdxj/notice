package neoproxy

import (
	"fmt"
	"time"
)

type News struct {
	Title      string
	UpdateTime time.Time
	Content    string
}

// Title: xxx
// Update time: xxx
// Content: xxx
func (n *News) String() string {
	msg := fmt.Sprintf("Title: %s\nUpdate time: %s\nContent: %s\n----------",
		n.Title, n.UpdateTime.Format("2006-01-02"), n.Content)
	return msg
}
