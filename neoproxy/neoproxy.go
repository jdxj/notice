package neoproxy

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/jdxj/notice/client"
	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/email"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
)

func NewFlow() *Flow {
	c := client.NewClientCookie(
		neoProxyCfg.Host,
		neoProxyCfg.Cookies,
		neoProxyCfg.Domain,
	)

	flow := &Flow{
		client:      c,
		dosageMutex: &sync.Mutex{},
		newsMutex:   &sync.Mutex{},
	}
	return flow
}

const (
	// 需要配合 neoProxyCfg.Host 使用
	LoginPage  = "/user"
	DosagePage = "/services"
	NewsPage   = "/news"
)

var (
	ErrLoginFailed = errors.New("login failed")

	neoProxyCfg *config.Neo
)

func init() {
	cfg, err := config.GetNeo()
	if err != nil {
		// todo: 在 init 之前检查配置正确性
		//panic(err)
	}
	neoProxyCfg = cfg
}

type Flow struct {
	client *http.Client

	// dosageMutex 保护以下字段
	dosageMutex *sync.Mutex
	lables      []string
	sample      []float64
	stat        map[string]float64

	// newsMutex 保护以下字段
	newsMutex *sync.Mutex
	news      *News
	newsTitle string
}

func (flow *Flow) VerifyLogin() error {
	resp, err := flow.client.Get(neoProxyCfg.Host + LoginPage)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
	}
	selection := doc.Find(".__cf_email__")

	cipherAddr, exist := selection.Attr("data-cfemail")
	if !exist {
		return fmt.Errorf("not found data-cfemail attr")
	}

	addr, err := decodeEmail(cipherAddr)
	if err != nil {
		return err
	}

	if addr != neoProxyCfg.User {
		return fmt.Errorf("verify email failed")
	}
	return nil
}

func (flow *Flow) UpdateDosage() {
	dosageURL := neoProxyCfg.Host + DosagePage + neoProxyCfg.Services
	resp, err := flow.client.Get(dosageURL)
	if err != nil {
		logs.Error("update dosage failed: %s", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("new doc failed: %s", err)
		return
	}

	selection := doc.Find("script").Last()
	script := selection.Text()
	//fmt.Printf("script: %s\n", script)

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
			logs.Error("parse dosage failed: %s", transfers[i])
		}
		stat[k] = dosage
	}

	// 保存采样
	flow.dosageMutex.Lock()
	flow.lables = lables
	flow.stat = stat
	flow.dosageMutex.Unlock()

}

func (flow *Flow) SendDosage() {
	subject := "Neo Proxy 用量统计"
	content := fmt.Sprintf("%s\n%s\n",
		flow.TodayUsed(),
		flow.TotalUsed())

	if err := email.SendSelf(subject, content); err != nil {
		logs.Error("send dosage failed: %s", err)
	}
}

func (flow *Flow) TotalUsed() string {
	var used float64

	flow.dosageMutex.Lock()
	for _, dosage := range flow.stat {
		used += dosage
	}
	flow.dosageMutex.Unlock()

	return fmt.Sprintf("总共用量: %.2fG", used/1024)
}

func (flow *Flow) TodayUsed() string {
	today := time.Now().Format("Jan 02")
	var dosage float64

	flow.dosageMutex.Lock()
	dosage = flow.stat[today]
	flow.dosageMutex.Unlock()

	return fmt.Sprintf("今日用量: %.2fM", dosage)
}

func (flow *Flow) SendLastNews() {
	var news *News
	flow.newsMutex.Lock()
	news = flow.news
	flow.newsMutex.Unlock()

	if news == nil {
		return
	}
	// 仍是之前的消息
	if news.Title == flow.newsTitle {
		return
	}

	subject := "新消息"
	content := news.String()
	if err := email.SendSelf(subject, content); err != nil {
		logs.Error("send last news failed: %s", err)
		return
	}

	flow.newsTitle = news.Title
}

func (flow *Flow) CrawlLastNews() {
	newsURLs, err := flow.crawlNewsList()
	if err != nil {
		logs.Error("crawl last news failed: %s", err)
		return
	}

	if len(newsURLs) <= 0 {
		logs.Warn("no news")
		return
	}
	lastNewsURL := newsURLs[0]
	resp, err := flow.client.Get(lastNewsURL)
	if err != nil {
		logs.Error("get failed: %s", err)
		return
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logs.Error("new doc failed: %s", err)
		return
	}

	selection := doc.Find(".card-body")
	// 解析 title
	title := selection.Find("h3").Text()
	// 解析 update time
	dateStr := selection.Find("small").Text()
	dateStrIdx := strings.Index(dateStr, "20")
	if dateStrIdx < 0 {
		logs.Error("not found news's date: %s", lastNewsURL)
		return
	}

	dateStr = dateStr[dateStrIdx:]
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		logs.Error("parse time failed: %s", err)
		return
	}
	// 解析 content
	content := selection.Find("p").Text()
	news := &News{
		Title:      title,
		UpdateTime: date,
		Content:    content,
	}

	flow.newsMutex.Lock()
	flow.news = news
	flow.newsMutex.Unlock()
}

func (flow *Flow) crawlNewsList() ([]string, error) {
	resp, err := flow.client.Get(neoProxyCfg.Host + NewsPage)
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

		newsURL = neoProxyCfg.Host + newsURL
		newsURLs = append(newsURLs, newsURL)
	})

	return newsURLs, nil
}

func decodeEmail(cipher string) (string, error) {
	if len(cipher) < 2 {
		return "", fmt.Errorf("cipher len not enough")
	}

	var address string
	key, err := strconv.ParseInt(cipher[:2], 16, 64)
	if err != nil {
		return "", err
	}

	for i := 2; i < len(cipher); i += 2 {
		x, err := strconv.ParseInt(cipher[i:i+2], 16, 64)
		if err != nil {
			logs.Error(err)
			return "", err
		}
		address += string(x ^ key)
	}
	return address, nil
}
