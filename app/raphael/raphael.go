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
	RSSExRomURL = "https://sourceforge.net/projects/evolution-x/rss?path=/raphael"
	RSSImURL    = "https://sourceforge.net/projects/unofficialbuilds/rss?path=/raphael/kernel"
)

func NewRaphael() *Raphael {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}

	r := &Raphael{
		client:  c,
		mutex:   &sync.Mutex{},
		mutexIm: &sync.Mutex{},
	}
	return r
}

type Raphael struct {
	client *http.Client

	// mutex 保护以下字段
	mutex *sync.Mutex
	item  *Item
	title string

	// mutexIm 保护以下字段
	mutexIm *sync.Mutex
	itemIm  *Item
	titleIm string
}

func (r *Raphael) UpdateItem() {
	data, err := r.readRespBody(RSSExRomURL)
	if err != nil {
		logs.Error("read all failed, err: %s, url: %s", err, RSSExRomURL)
		return
	}

	item, err := unmarshalRSS(data)
	if err != nil {
		logs.Error("unmarshal item failed, err: %s, url: %s", err, RSSExRomURL)
		return
	}

	r.mutex.Lock()
	r.item = item
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

func (r *Raphael) UpdateItemIm() {
	data, err := r.readRespBody(RSSImURL)
	if err != nil {
		logs.Error("read all failed, err: %s, url: %s", err, RSSImURL)
		return
	}

	item, err := unmarshalRSS(data)
	if err != nil {
		logs.Error("unmarshal item failed, err: %s, url: %s", err, RSSImURL)
		return
	}

	r.mutexIm.Lock()
	r.itemIm = item
	r.mutexIm.Unlock()
}

func (r *Raphael) readRespBody(url string) ([]byte, error) {
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

func (r *Raphael) SendUpdateIm() {
	var item *Item
	r.mutexIm.Lock()
	item = r.itemIm
	r.mutexIm.Unlock()
	// 没初始化
	if item == nil {
		return
	}
	// 没有更新
	if item.Title.Data == r.titleIm {
		return
	}
	r.titleIm = item.Title.Data

	subject := "iMMENSITY 已更新"
	content, _ := xml.MarshalIndent(item, "", "    ")
	if err := email.SendSelfBytes(subject, content); err != nil {
		logs.Error("send updateIm failed: %s", err)
	}
}
