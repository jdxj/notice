package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func main() {
	cookieStr := `__cfduid=d4e553e43f7c332bc446c52508fe514ce1577152713; an_visitation_log=6cd8485ba55c49e2826b888296aa90c9; _fbp=fb.1.1577152718125.1425761279; _ga=GA1.2.1523528149.1577152718; track_id=028e6a2856044eb5ab4dcf4ce209bd47; _gid=GA1.2.823441203.1578210295; _gat_gtag_UA_115741053_1=1; _agentneo_session=b2VEaGNxNnc2cUhwRlRMUkxPZnhQbGhURE5Ea3AyZDlONElsVDRNL3RubitHTEJHRFp4cmkxN2tZcGRrNC9YdFltNThoakxKYWtiNW9CZ1J2RlV4SUh0a2syZWZML1lCRExDNzQwMU9Vbjkyc2lsTWxqOUVKeUxYdnlCQnpvL09VZ0FGZ0RHdnFrTS90Mm9velkvRzdmOXlCb3NWYUxmSjhtcW9zMXdybHJ2eVE0ZHVYdVdRRGFNVXpaYVE1OXRKRFYvZUR2dDQ3SzhmeFlYQ0VpT3A4RXozNFEzcDU0ZlNlbjlFWml4YWVXL3pmWFQySGpHNEJMRU4ya2g4b1FxY2dJK1pIcGFPWDdKNy93UHlvQ1JQU0JkMnVueXlIaDhzeC9iNFNkbENMRXV2V01TZjJWT0Fxa3A1ekVGQ09rYzVNclBwUGFILzNtcjlBcENUYjN2VFVnPT0tLXQ3eGxna2Y3ZkRESGJqQ1ZybENyemc9PQ%3D%3D--42bf98d28696a1aaae71ca81765177809911fdc6`
	service := "https://neoproxy.org/services/b9a8283c7fd7407bb9c13e04af5b0fad"

	cookies, err := StrToCookies(cookieStr, ".neoproxy.org")
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	jar, err := cookiejar.New(nil)
	if err != nil {
		panic(err)
	}

	cookieURL, err := url.Parse("https://neoproxy.org")
	if err != nil {
		panic(err)
	}

	jar.SetCookies(cookieURL, cookies)

	client.Jar = jar

	resp, err := client.Get(service)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", data)

}


func StrToCookies(cookiesStr, domain string) ([]*http.Cookie, error) {
	if domain == "" {
		return nil, fmt.Errorf("invaild domain")
	}

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
			//log.MyLogger.Warn("%s not found '=' in cookie part: %s", log.Log_Log, part)
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

	if len(cookies) == 0 {
		return nil, fmt.Errorf("invalid cookie")
	}
	return cookies, nil
}
