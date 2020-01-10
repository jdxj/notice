package neoproxy

import (
	"fmt"
	"net/http"
	"notice/module"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
)

func NewFlow(npc *module.NeoProxyConfig, ec *module.EmailConfig) (*Flow, error) {
	if npc == nil {
		return nil, fmt.Errorf("invalid NeoProxyConfig")
	}

	es, err := module.NewEmailSender(ec)
	if err != nil {
		return nil, err
	}

	flow := &Flow{
		config:      npc,
		stop:        make(chan struct{}),
		wg:          &sync.WaitGroup{},
		emailSender: es,
	}

	flow.client, err = module.NewHTTPClientWithCookie(npc.URL, npc.Cookie, npc.Domain)
	if err != nil {
		return nil, err
	}
	return flow, nil
}

type Flow struct {
	config      *module.NeoProxyConfig
	client      *http.Client
	emailSender *module.EmailSender

	lables []string
	stat   map[string]float64

	stop chan struct{}
	wg   *sync.WaitGroup
}

func (flow *Flow) Start() {
	wg := flow.wg

	wg.Add(1)
	go func() {
		flow.NotifyAt11pm()
		wg.Done()
	}()

	logs.Info("flow started")
}

func (flow *Flow) Stop() {
	close(flow.stop)
	flow.wg.Wait()

	logs.Info("flow stopped")
}

func (flow *Flow) NotifyAt11pm() {
	// 先运行一次
	flow.update()
	flow.sendNotice()

	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	for {
		now := time.Now()
		tomorrow11pm := time.Unix((now.Unix()/86400+1)*86400, 0).Add(15 * time.Hour)
		dur := tomorrow11pm.Sub(now)
		timer.Reset(dur)

		select {
		case <-flow.stop:
			logs.Info("stop notify at 11pm")
			return
		case <-timer.C:
			logs.Debug("notify at 11pm")
			flow.update()
			flow.sendNotice()
		}
	}
}

func (flow *Flow) update() {
	resp, err := flow.client.Get(flow.config.Service)
	if err != nil {
		logs.Error("error when update flow: %s", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("error when new doc")
		return
	}

	selection := doc.Find("script").Last()
	script := selection.Text()

	// 解析 lables
	lablesStart := strings.Index(script, `[`)
	lablesEnd := strings.Index(script, `]`)
	lablesStr := script[lablesStart+1 : lablesEnd]
	lablesStr = strings.ReplaceAll(lablesStr, `'`, ``)
	lables := strings.Split(lablesStr, `,`)

	// 解析 transfers
	transfersStart := strings.LastIndex(script, `[`)
	transfersEnd := strings.LastIndex(script, `]`)
	transfersStr := script[transfersStart+1 : transfersEnd]
	transfers := strings.Split(transfersStr, `,`)

	if len(lables) != len(transfers) || len(lables) <= 0 {
		logs.Error("daily data invalid")
		return
	}

	stat := make(map[string]float64)
	for i, k := range lables {
		dosage, err := strconv.ParseFloat(transfers[i], 64)
		if err != nil {
			logs.Error("error when parse dosage: %s", transfers[i])
		}
		stat[k] = dosage
	}

	flow.lables = lables
	flow.stat = stat
	logs.Debug("update flow success")

}

func (flow *Flow) TotalUsed() string {
	var used float64

	for _, dosage := range flow.stat {
		used += dosage
	}
	return fmt.Sprintf("%.2fG", used/1024)
}

func (flow *Flow) TodayUsed() string {
	today := time.Now().Format("Jan 02")
	return fmt.Sprintf("%.2fM", flow.stat[today])
}

func (flow *Flow) sendNotice() {
	subject := "用量统计"
	content := fmt.Sprintf("当日用量: %s, 总用量: %s", flow.TodayUsed(), flow.TotalUsed())

	e := flow.emailSender
	e.SendMsg(subject, content)
}
