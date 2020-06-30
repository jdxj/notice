package sourceforge

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/astaxie/beego/logs"

	"github.com/jdxj/notice/client"
	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/email"
)

func NewSourceforge(rssURL string) *Sourceforge {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}

	r := &Sourceforge{
		url:    rssURL,
		client: c,
		mutex:  &sync.Mutex{},
	}
	return r
}

type Sourceforge struct {
	url string

	client *http.Client

	// mutex 保护以下字段
	mutex *sync.Mutex
	item  *Item
	title string
}

func (r *Sourceforge) UpdateItem() {
	req, err := client.NewRequestUserAgentGet(r.url)
	if err != nil {
		logs.Error("new request failed: %s", err)
		return
	}
	data, err := client.ReadResponseBody(r.client, req)
	if err != nil {
		logs.Error("read response body failed: %s", err)
		return
	}

	item, err := unmarshalRSS(data)
	if err != nil {
		logs.Error("unmarshal item failed, err: %s, url: %s", err, r.url)
		return
	}

	r.mutex.Lock()
	r.item = item
	r.mutex.Unlock()
}

func (r *Sourceforge) SendUpdate() {
	var item *Item
	r.mutex.Lock()
	item = r.item
	r.mutex.Unlock()
	// 没初始化
	if item == nil {
		return
	}
	// 没有更新
	if item.Title.Data == r.title {
		return
	}
	r.title = item.Title.Data

	subject := fmt.Sprintf("<%s> 已更新", item.Title.Data)
	content, _ := xml.MarshalIndent(item, "", "    ")

	ds := config.DataStorage
	emailCfg, err := ds.GetEmail()
	if err != nil {
		logs.Error("get sourceforge config failed: %s", err)
		return
	}

	if err := email.SendTextBytes(subject, content, emailCfg.Addr); err != nil {
		logs.Error("send text by bytes failed, subject: %s, content: %s",
			subject, content)
	}
}
