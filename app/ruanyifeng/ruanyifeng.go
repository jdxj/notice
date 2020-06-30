package ruanyifeng

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/email"

	"github.com/astaxie/beego/logs"
	"github.com/jdxj/notice/client"
)

const (
	RSSURL = "http://www.ruanyifeng.com/blog/atom.xml"
)

func NewRuanYiFeng() *RuanYiFeng {
	jar, _ := cookiejar.New(nil)
	c := &http.Client{
		Jar: jar,
	}

	ryf := &RuanYiFeng{
		client:     c,
		entryMutex: &sync.Mutex{},
	}
	return ryf
}

type RuanYiFeng struct {
	client *http.Client

	// entryMutex 保护以下字段
	entryMutex *sync.Mutex
	entry      *Entry
	title      string
}

func (ryf *RuanYiFeng) UpdateEntry() {
	req, err := client.NewRequestUserAgentGet(RSSURL)
	if err != nil {
		logs.Error("ruanyifeng new req failed: %s", err)
		return
	}

	data, err := client.ReadResponseBody(ryf.client, req)
	if err != nil {
		logs.Error("ruanyifeng read body failed: %s", err)
		return
	}

	entry, err := unmarshalFeed(data)
	if err != nil {
		logs.Error("unmarshal entry failed, err: %s, url: %s", err, RSSURL)
		return
	}

	ryf.entryMutex.Lock()
	ryf.entry = entry
	ryf.entryMutex.Unlock()
}

func (ryf *RuanYiFeng) SendUpdate() {
	var entry *Entry
	ryf.entryMutex.Lock()
	entry = ryf.entry
	ryf.entryMutex.Unlock()

	if entry == nil {
		return
	}
	if entry.Title == ryf.title {
		return
	}
	ryf.title = entry.Title

	subject := fmt.Sprintf("<阮一峰的网络日志> 已更新")
	content := entry.Content.Data

	ryfCfg, err := config.GetRYF()
	if err != nil {
		logs.Error("get ryf config failed: %s", err)
		return
	}
	if err := email.SendHTMLBytes(subject, content, ryfCfg.Users...); err != nil {
		logs.Error("send html by bytes failed, subject: %s, content: %s",
			subject, content)
	}
}
