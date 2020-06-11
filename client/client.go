package client

import (
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
)

const (
	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36"
)

func NewClientCookie(rawURL, cookiesStr, domain string) *http.Client {
	cookies := StringToCookies(cookiesStr, domain)
	URL, _ := url.Parse(rawURL)

	jar, _ := cookiejar.New(nil)
	jar.SetCookies(URL, cookies)

	client := &http.Client{
		Jar: jar,
	}
	return client
}

func NewRequestUserAgent(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", UserAgent)
	return req, nil
}

func StringToCookies(cookiesStr, domain string) []*http.Cookie {
	// 过滤引号
	cookiesStr = strings.ReplaceAll(cookiesStr, `"`, ``)
	// 过滤空格
	cookiesStr = strings.ReplaceAll(cookiesStr, ` `, ``)
	// 划分
	cookiesParts := strings.Split(cookiesStr, ";")

	var cookies []*http.Cookie
	for _, part := range cookiesParts {
		idx := strings.Index(part, "=")
		if idx < 0 {
			logs.Warn("not found '=', domain: %s, part: %'", domain, part)
			continue
		}
		k := part[:idx]
		v := part[idx+1:]

		cookie := &http.Cookie{
			Name:     k,
			Value:    v,
			Path:     "/",
			Domain:   domain,
			Expires:  time.Now().Add(time.Hour * 24 * 365), // 一年后过期
			Secure:   false,
			HttpOnly: false,
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}
