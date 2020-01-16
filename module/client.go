package module

import (
	"io"
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

func NewHTTPRequestWithUserAgent(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/79.0.3945.117 Safari/537.36")
	return req, nil
}
