package module

import (
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"notice/utils"
)

func NewHTTPClientWithCookie(rawURL, cookiesStr, domain string) (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	cookies, err := utils.StrToCookies(cookiesStr, domain)
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	jar.SetCookies(u, cookies)

	client := &http.Client{}
	client.Jar = jar
	return client, nil
}
