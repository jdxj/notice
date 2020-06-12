package raphael

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/jdxj/notice/email"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/client"
)

const (
	RSSURL = "https://sourceforge.net/projects/evolution-x/rss?path=/raphael"
)

func NewRaphael() *Raphael {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}

	r := &Raphael{
		client: c,
		mutex:  &sync.Mutex{},
	}
	return r
}

type Raphael struct {
	client *http.Client

	// mutex 保护以下字段
	mutex *sync.Mutex
	item  *Item
	title string
}

func (r *Raphael) UpdateItem() {
	req, err := client.NewRequestUserAgent(http.MethodGet, RSSURL, nil)
	if err != nil {
		logs.Error("new req user agent failed: %s", err)
		return
	}
	resp, err := r.client.Do(req)
	if err != nil {
		logs.Error("do req failed: %s", err)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("read all failed: %s", err)
		return
	}

	rss := &RSS{}
	if err := xml.Unmarshal(data, rss); err != nil {
		logs.Error("unmarshal xml failed: %s", err)
		return
	}
	if rss.Channel == nil || len(rss.Channel.Items) <= 0 {
		logs.Warn("not found raphael update")
		return
	}

	r.mutex.Lock()
	r.item = rss.Channel.Items[0]
	r.mutex.Unlock()
}

func (r *Raphael) SendUpdate() {
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

	subject := "Evolution-x Raphael 已更新"
	content, _ := xml.MarshalIndent(item, "", "    ")
	if err := email.SendSelfBytes(subject, content); err != nil {
		logs.Error("send update failed: %s", err)
	}
}
