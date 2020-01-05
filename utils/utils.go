package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

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
