package neoproxy

import (
	"encoding/json"
	"fmt"
	"notice/module"
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

func TestCrawlNews(t *testing.T) {
	cfg, err := module.ReadConfig()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	flow, err := NewFlow(cfg.NeoProxy, cfg.Email)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	newsURLs, err := flow.crawlNewsList()
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	fmt.Printf("len: %d\n", len(newsURLs))
	for _, v := range newsURLs {
		fmt.Println(v)
	}

	news, err := flow.crawlLastNews()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s\n", news)
}

func TestFlow_VerifyLogin(t *testing.T) {
	cfg, err := module.ReadConfig()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	flow, err := NewFlow(cfg.NeoProxy, cfg.Email)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	if err := flow.VerifyLogin(); err != nil {
		t.Fatalf("%s", err)
	}
}

func TestTimeEqual(t *testing.T) {
	now := time.Now()
	if now != now {
		t.Fatalf("%s", "time can not compare")
	}
}

func TestDecodeEmail(t *testing.T) {
	address, err := decodeEmail("a990919c9e9c909b9f9be9d8d887cac6c4")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	if address != "985759262@qq.com" {
		t.Fatalf("decode failed\n")
	}
}

func TestFlow_SerializeCookie(t *testing.T) {
	cfg, err := module.ReadConfig()
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	flow, err := NewFlow(cfg.NeoProxy, cfg.Email)
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	cookieStr, err := flow.serializeCookie()
	if err != nil {
		t.Fatalf("%s", err)
	}
	fmt.Println(cookieStr)
}

func TestJsonMarshal(t *testing.T) {
	cfg, err := module.ReadConfig()
	if err != nil {
		t.Fatalf("%s\n", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	fmt.Printf("%s", data)
}
