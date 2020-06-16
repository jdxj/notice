package sourceforge

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/jdxj/notice/config"

	"github.com/jdxj/notice/email"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/client"
)

func NewSourceforge(rssURL string, emailCfg *config.Email) *Sourceforge {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}

	r := &Sourceforge{
		url:    rssURL,
		client: c,
		sender: email.NewSender(emailCfg),
		mutex:  &sync.Mutex{},
	}
	return r
}

type Sourceforge struct {
	url string

	client *http.Client
	sender *email.Sender

	// mutex 保护以下字段
	mutex *sync.Mutex
	item  *Item
	title string
}

func (r *Sourceforge) UpdateItem() {
	data, err := r.readRespBody(r.url)
	if err != nil {
		logs.Error("read all failed, err: %s, url: %s", err, r.url)
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

	sender := r.sender
	if err := sender.SendTextSelfBytes(subject, content); err != nil {
		logs.Error("send update failed: %s", err)
		return
	}
}

func (r *Sourceforge) readRespBody(url string) ([]byte, error) {
	req, err := client.NewRequestUserAgent(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}
