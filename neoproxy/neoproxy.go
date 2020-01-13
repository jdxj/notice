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

const (
	CollectingFrequency  = time.Minute
	SamplesLimit         = 1440 // 24h * 60m
	NotificationInterval = 10 * time.Minute
	DailyDosageLimit     = 600
)

const (
	HomePage = "https://neoproxy.org"
	NewsPage = "https://neoproxy.org/news"
)

type Flow struct {
	config      *module.NeoProxyConfig
	client      *http.Client
	emailSender *module.EmailSender

	lables       []string
	stat         map[string]float64
	sample       []float64
	lastNewsDate time.Time

	stop chan struct{}
	wg   *sync.WaitGroup
}

func (flow *Flow) Start() {
	// 这里使用通知的目的是: 如果程序由 systemd 重启过,
	// 说明发生了异常终止, 所以需要通知.
	subject := "notice start"
	if err := flow.emailSender.SendMsgString(subject, time.Now().Format(time.RFC1123)); err != nil {
		logs.Error("%s", err)
	}

	wg := flow.wg

	wg.Add(1)
	go func() {
		flow.CollectingSamples()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		flow.NotifyAt11pm()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		flow.NotifyExceedDailyDosageLimit()
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		flow.NotifyLastNews()
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
	subject := "用量统计"

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
			content := fmt.Sprintf("当日用量: %s, 总用量: %s", flow.TodayUsed(), flow.TotalUsed())
			err := flow.emailSender.SendMsgString(subject, content)
			if err != nil {
				logs.Error("%s", err)
			}
		}
	}
}

// NotifyExceedDailyDosageLimit 用于当超过每日平均用量时发出通知.
// 如果当天超过用量限制, 那么只通知一次.
func (flow *Flow) NotifyExceedDailyDosageLimit() {
	var alreadyNotified bool
	var notificationTime time.Time
	subject := fmt.Sprintf("超过每日用量限制!")

	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-flow.stop:
			logs.Info("stop notify daily dosage limit")
			return

		case <-ticker.C:
			// 新的一天到了
			if time.Now().Unix()/86400 != notificationTime.Unix()/86400 {
				alreadyNotified = false
			}

			if flow.todayUsed() > DailyDosageLimit && !alreadyNotified {
				content := fmt.Sprintf("当日用量: %s, 总用量: %s", flow.TodayUsed(), flow.TotalUsed())
				if err := flow.emailSender.SendMsgString(subject, content); err != nil {
					logs.Error("error when send exceed daily dosage notice: %s", err)
				} else {
					alreadyNotified = true
					notificationTime = time.Now()
				}
			}
		}
	}
}

// CollectingSamples 用于收集用量数据.
// 所有的通知方法不主动请求用量数据, 它们统计的数据都由 CollectingSamples
// 同一收集. 应该保证 CollectingSamples 采集的频率有一个合适的值.
func (flow *Flow) CollectingSamples() {
	ticker := time.NewTicker(CollectingFrequency)
	defer ticker.Stop()

	for {
		select {
		case <-flow.stop:
			logs.Info("stop collecting samples")
			return

		case <-ticker.C:
			flow.update()
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

	// 保存数据
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

	// 保存采样
	today := time.Now().Format("Jan 02")
	dosage := flow.stat[today]
	flow.sample = append(flow.sample, dosage)
	if len(flow.sample) > SamplesLimit {
		flow.sample = flow.sample[1:]
	}
}

func (flow *Flow) TotalUsed() string {
	var used float64

	for _, dosage := range flow.stat {
		used += dosage
	}
	return fmt.Sprintf("%.2fG", used/1024)
}

func (flow *Flow) TodayUsed() string {
	return fmt.Sprintf("%.2fM", flow.todayUsed())
}

func (flow *Flow) todayUsed() float64 {
	today := time.Now().Format("Jan 02")
	return flow.stat[today]
}

func (flow *Flow) NotifyLastNews() {
	subject := "新的消息"

	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for {
		select {
		case <-flow.stop:
			logs.Info("stop notify last news")
			return

		case <-ticker.C:
			news, err := flow.crawlLastNews()
			if err != nil {
				logs.Error("%s", err)
				continue
			}
			if flow.lastNewsDate == news.UpdateTime {
				continue
			}
			flow.lastNewsDate = news.UpdateTime

			err = flow.emailSender.SendMsgString(subject, news.String())
			if err != nil {
				logs.Error("%s", err)
			}
		}
	}
}

func (flow *Flow) crawlLastNews() (*News, error) {
	newsURLs, err := flow.crawlNewsList()
	if err != nil {
		return nil, err
	}
	if len(newsURLs) <= 0 {
		return nil, fmt.Errorf("no news")
	}

	lastNewsURL := newsURLs[0]
	resp, err := flow.client.Get(lastNewsURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	selection := doc.Find(".card-body")

	// 解析 title
	title := selection.Find("h3").Text()

	// 解析 update time
	dateStr := selection.Find("small").Text()
	dateStrIdx := strings.Index(dateStr, "20")
	if dateStrIdx < 0 {
		return nil, fmt.Errorf("not found news's date: %s", lastNewsURL)
	}
	dateStr = dateStr[dateStrIdx:]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, err
	}

	// 解析 content
	content := selection.Find("p").Text()

	news := &News{
		Title:      title,
		UpdateTime: date,
		Content:    content,
	}
	return news, nil
}

func (flow *Flow) crawlNewsList() ([]string, error) {
	resp, err := flow.client.Get(NewsPage)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	selection := doc.Find(".list-group-item-action")

	var newsURLs []string
	selection.Each(func(i int, sel *goquery.Selection) {
		newsURL, exist := sel.Attr("href")
		if !exist {
			return
		}

		newsURL = HomePage + newsURL
		newsURLs = append(newsURLs, newsURL)
	})

	return newsURLs, nil
}
