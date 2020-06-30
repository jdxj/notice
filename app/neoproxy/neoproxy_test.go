package neoproxy

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestTimeFormat(t *testing.T) {
	str := time.Now().Add(time.Hour * 24 * 30).Format("Jan 02")
	fmt.Println(str)
}

func TestWaitGroup(t *testing.T) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		time.Sleep(5 * time.Second)
		wg.Done()
	}()

	//wg.Wait()
}

func TestTimer(t *testing.T) {
	now := time.Now()
	//tomorrow11pm := time.Unix((now.Unix()/86400+1)*86400, 0).Add(23 * time.Hour)
	tomorrow11pm := time.Unix((now.Unix()/86400+1)*86400, 0).Add(15 * time.Hour)
	fmt.Println(tomorrow11pm)
}

func TestNews_String(t *testing.T) {
	news := &News{
		Title:      "test",
		UpdateTime: time.Now(),
		Content:    "abc def",
	}
	fmt.Println(news)

	news = &News{
		Title:      "test",
		UpdateTime: time.Now(),
		Content:    "abc def",
	}
	fmt.Println(news)

	news = &News{
		Title:      "test",
		UpdateTime: time.Now(),
		Content:    "abc def",
	}
	fmt.Println(news)
}

func TestFlow_VerifyLogin(t *testing.T) {
	flow := NewFlow()
	flow.NotifyDosage()
}
