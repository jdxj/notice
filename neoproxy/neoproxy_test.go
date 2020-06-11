package neoproxy

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/jdxj/notice/config"
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
	if err := flow.VerifyLogin(); err != nil {
		t.Fatalf("%s\n", err)
	}

	flow.UpdateDosage()
	flow.SendDosage()
}

func TestFlow_CrawlLastNews(t *testing.T) {
	flow := NewFlow()

	if err := flow.VerifyLogin(); err != nil {
		t.Fatalf("%s\n", err)
	}

	flow.CrawlLastNews()
}

func TestNeoConfig(t *testing.T) {
	neo := &config.Neo{
		Host:     "https://neogate.co",
		Domain:   ".neogate.co",
		Cookies:  "",
		Services: "",
		User:     "985759262@qq.com",
	}
	if err := config.SetNeo(neo); err != nil {
		t.Fatalf("%s\n", err)
	}
	config.Close()
	//fmt.Printf("%v", *neoProxyCfg)
}

func TestTimeEqual(t *testing.T) {
	now := time.Now()
	if now != now {
		t.Fatalf("%s", "time can not compare")
	}
}

func TestDecodeEmail(t *testing.T) {
	//address, err := decodeEmail("a990919c9e9c909b9f9be9d8d887cac6c4")
	address, err := decodeEmail("d1e8e9e4e6e4e8e3e7e391a0a0ffb2bebc")
	if err != nil {
		t.Fatalf("%s\n", err)
	}
	if address != "985759262@qq.com" {
		t.Fatalf("decode failed\n")
	}
}

func TestFlow_SendLastNews(t *testing.T) {
	flow := NewFlow()
	if err := flow.VerifyLogin(); err != nil {
		t.Fatalf("%s\n", err)
	}

	flow.CrawlLastNews()
	flow.SendLastNews()
}
