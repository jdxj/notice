package bilibili

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/jdxj/notice/client"
	"github.com/jdxj/notice/config"
	"github.com/jdxj/notice/email"

	"github.com/astaxie/beego/logs"
)

const (
	Domain   = ".bilibili.com"
	LoginAPI = "https://api.bilibili.com/x/web-interface/nav"
	SignAPI  = "https://api.bilibili.com/x/member/web/coin/log?jsonp=jsonp"
	CoinAPI  = "https://account.bilibili.com/site/getCoin"
	WebSite  = "https://www.bilibili.com"

	TimeLayout = "2006-01-02 15:04:05"
)

var (
	ErrSigRecNotFound = fmt.Errorf("sign-in record not found")
)

func NewBiliBili() *BiliBili {
	webURL, _ := url.Parse(WebSite)
	bili := &BiliBili{
		client: &http.Client{},
		url:    webURL,
	}
	return bili
}

type BiliBili struct {
	client *http.Client
	url    *url.URL
}

func (bili *BiliBili) KeepCollect() {
	biliCfg, err := config.GetBiliBili()
	if err != nil {
		logs.Error("get bilibili config failed: %s", err)
		return
	}

	user := make(map[string]*cookiejar.Jar)
	for addr, cookies := range biliCfg.Cookies {
		jar := client.NewJarCookie(bili.url, cookies, Domain)
		user[addr] = jar
	}
	for addr, jar := range user {
		bili.client.Jar = jar
		bili.collectCoins(addr)
	}
}

func (bili *BiliBili) collectCoins(addr string) {
	// 1.
	if err := bili.verifyLogin(); err != nil {
		logs.Error("verify login failed: %s", err)
		return
	}
	// 2.
	if err := bili.sign(); err != nil {
		logs.Error("sign failed: %s", err)
		return
	}
	// 3.
	if err := bili.verifyMultiCheckSign(); err != nil {
		if err != ErrSigRecNotFound {
			logs.Error("verify multi check sing failed: %s", err)
			return
		}

		subject := "<BiliBili Coins Number>"
		content := fmt.Sprintf("没有查询到最新硬币余量")
		if err := email.SendText(subject, content, addr); err != nil {
			logs.Error("send email failed, subject: %s, content: %s",
				subject, content)
		}
		return
	}
	// 4.
	if err := bili.sendCoinsNum(addr); err != nil {
		logs.Error("%s", err)
	}
}

// 1. 验证登陆
func (bili *BiliBili) verifyLogin() error {
	req, _ := client.NewRequestUserAgentGet(LoginAPI)
	data, err := client.ReadResponseBody(bili.client, req)
	if err != nil {
		return err
	}

	data, err = unmarshalAPIResponse(data)
	if err != nil {
		return err
	}
	loginInfo := &LoginInfo{}
	if err := json.Unmarshal(data, loginInfo); err != nil {
		return err
	}
	if !loginInfo.IsLogin {
		return fmt.Errorf("login info: %v", loginInfo.IsLogin)
	}
	return nil
}

// 2. 签到
func (bili *BiliBili) sign() error {
	req, _ := client.NewRequestUserAgentGet(WebSite)
	resp, err := bili.client.Do(req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// 3. 检查签到结果
func (bili *BiliBili) checkSign() error {
	req, _ := client.NewRequestUserAgentGet(SignAPI)
	data, err := client.ReadResponseBody(bili.client, req)
	if err != nil {
		return err
	}

	data, err = unmarshalAPIResponse(data)
	if err != nil {
		return err
	}
	signInfo := &SignInfo{}
	if err := json.Unmarshal(data, signInfo); err != nil {
		return err
	}

	signEntry := signInfo.List[0]
	curDate, _ := time.Parse(TimeLayout, signEntry.Time)
	now := time.Now()

	if curDate.Year() != now.Year() &&
		curDate.Month() != now.Month() &&
		curDate.Day() != now.Day() {

		return ErrSigRecNotFound
	}
	return nil
}

// 3.a 多次验证是否签到成功
func (bili *BiliBili) verifyMultiCheckSign() error {
	// 1+2+4+8+...
	checkNum := 6 // 检查次数
	timer := time.NewTimer(time.Minute)
	defer timer.Stop()

	for i := 0; i < checkNum; i++ {
		dur := math.Pow(2, float64(i))
		timer.Reset(time.Duration(dur) * time.Second)
		<-timer.C

		err := bili.checkSign()
		if err == nil || err != ErrSigRecNotFound {
			return err
		}
		// 如果 err == ErrSigRecNotFound 则继续下一次循环
	}
	return ErrSigRecNotFound
}

// 4. 发送硬币数量
func (bili *BiliBili) sendCoinsNum(addr string) error {
	req, _ := client.NewRequestUserAgentGet(CoinAPI)
	data, err := client.ReadResponseBody(bili.client, req)
	if err != nil {
		return err
	}

	apiResp := &APIResponse{}
	if err := json.Unmarshal(data, apiResp); err != nil {
		return err
	}
	if !apiResp.Status {
		return fmt.Errorf("verify status failed, status: %v", apiResp.Status)
	}

	coinIfo := &CoinInfo{}
	if err := json.Unmarshal(apiResp.Data, coinIfo); err != nil {
		return err
	}

	subject := fmt.Sprintf("<BiliBili Coins Number>")
	content := fmt.Sprintf("硬币余额: %d", coinIfo.Money)
	return email.SendText(subject, content, addr)
}
