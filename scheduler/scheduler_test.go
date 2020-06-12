package scheduler

import (
	"fmt"
	"testing"
	"time"
)

func TestScheduler(t *testing.T) {
	_, err := scheduler.AddFunc("0 * * * *", func() {
		fmt.Printf("%s - hello world\n", time.Now())
	})
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	scheduler.Start()

	time.Sleep(10 * time.Minute)
	scheduler.Stop()
}
