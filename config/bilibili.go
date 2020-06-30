package config

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

var (
	bilibiliKey = []byte("bilibili")
)

type BiliBili struct {
	Cookies map[string]string `json:"cookies"`
}

func (bili *BiliBili) String() string {
	data, _ := json.MarshalIndent(bili, "", "    ")
	return fmt.Sprintf("%s", data)
}

func GetBiliBili() (*BiliBili, error) {
	data, err := get(bilibiliKey)
	if err != nil {
		return nil, err
	}

	bili := &BiliBili{}
	return bili, json.Unmarshal(data, bili)
}

func SetBiliCookie(emailAddr, cookie string) error {
	bili := &BiliBili{}
	bili.Cookies = map[string]string{
		emailAddr: cookie,
	}
	data, _ := json.Marshal(bili)
	return set(bilibiliKey, data)
}

func AddBiliCookie(emailAddr, cookie string) error {
	if emailAddr == "" || cookie == "" {
		return fmt.Errorf("email or cookie is empty")
	}

	bili, err := GetBiliBili()
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("add cookie failed: %s", err)
		}
		return SetBiliCookie(emailAddr, cookie)
	}

	bili.Cookies[emailAddr] = cookie
	data, _ := json.Marshal(bili)
	return set(bilibiliKey, data)
}
