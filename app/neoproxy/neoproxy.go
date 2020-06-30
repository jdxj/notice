package neoproxy

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/jdxj/notice/client"

	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/email"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego/logs"
)

func NewFlow() *Flow {
	flow := &Flow{
		c: &http.Client{},
	}
	return flow
}

const (
	// 需要配合 neoProxyCfg.Host 使用
	LoginPage  = "/user"
	DosagePage = "/services"
	NewsPage   = "/news"
)

type Flow struct {
	c      *http.Client
	lables []string
	stat   map[string]float64

	// newsMutex 保护以下字段
	//newsMutex *sync.Mutex
	//news      *News
	//newsTitle string
}

func (flow *Flow) NotifyDosage() {
	neoCfg, err := config.GetNeo()
	if err != nil {
		logs.Error("get neo config failed: %s", err)
		return
	}

	URL, err := url.Parse(neoCfg.Host)
	if err != nil {
		logs.Error("parse url failed: %s", err)
		return
	}

	flow.c.Jar = client.NewJarCookie(URL, neoCfg.Cookies, neoCfg.Domain)
	// 1.
	if err := flow.verifyLogin(neoCfg); err != nil {
		logs.Error("verify neo login failed: %s", err)
		return
	}
	// 2.
	if err := flow.updateDosage(neoCfg); err != nil {
		logs.Error("update dosage failed: %s", err)
		return
	}
	// 3.
	flow.sendDosage(neoCfg.User)
}

func (flow *Flow) verifyLogin(neoCfg *config.Neo) error {
	c := flow.c
	resp, err := c.Get(neoCfg.Host + LoginPage)
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

	if addr != neoCfg.User {
		return fmt.Errorf("verify email failed")
	}
	return nil
}

func (flow *Flow) updateDosage(neoCfg *config.Neo) error {
	c := flow.c
	dosageURL := neoCfg.Host + DosagePage + neoCfg.Services
	resp, err := c.Get(dosageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return err
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
		return fmt.Errorf("daily data invalid")
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

	flow.lables = lables
	flow.stat = stat
	return nil
}

func (flow *Flow) sendDosage(addr string) {
	subject := "Neo Proxy 用量统计"
	content := fmt.Sprintf("%s\n%s\n",
		flow.TodayUsed(),
		flow.TotalUsed())

	if err := email.SendText(subject, content, addr); err != nil {
		logs.Error("send email failed, subject: %s, content: %s",
			subject, content)
	}
}

func (flow *Flow) TotalUsed() string {
	var used float64
	for _, dosage := range flow.stat {
		used += dosage
	}
	return fmt.Sprintf("总共用量: %.2fG", used/1024)
}

func (flow *Flow) TodayUsed() string {
	today := time.Now().Format("Jan 02")
	dosage := flow.stat[today]
	return fmt.Sprintf("今日用量: %.2fM", dosage)
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
