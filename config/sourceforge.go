package config

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

var (
	sourceforgeKey = []byte("sourceforge")
)

type Sourceforge struct {
	SubsAddr []string `json:"subs_addr"`
}

func (c *Cache) GetSourceforge() (*Sourceforge, error) {
	data, err := c.Get(sourceforgeKey)
	if err != nil {
		return nil, fmt.Errorf("cat not get data, key: %s, err: %s",
			sourceforgeKey, err)
	}

	sf := &Sourceforge{}
	return sf, json.Unmarshal(data, sf)
}

func (c *Cache) SetSubsAddr(rssURL string) error {
	sf := &Sourceforge{}
	sf.SubsAddr = append(sf.SubsAddr, rssURL)
	data, _ := json.Marshal(sf)
	return c.Set(sourceforgeKey, data)
}

func (c *Cache) AddSubsAddr(rssURL string) error {
	// url 检查
	if rssURL == "" {
		return fmt.Errorf("rss address invalid")
	}

	// 这里要使用 badger 中定义的 error, 所以没有调用 Cache.GetSourceforge()
	data, err := c.Get(sourceforgeKey)
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("cat not get data, key: %s, err: %s",
				sourceforgeKey, err)
		}

		return c.SetSubsAddr(rssURL)
	}

	// 取出已缓存的 url
	sf := &Sourceforge{}
	if err := json.Unmarshal(data, sf); err != nil {
		return err
	}

	// 查重
	for _, v := range sf.SubsAddr {
		if v == rssURL {
			return fmt.Errorf("duplicate subscription address: %s", rssURL)
		}
	}

	sf.SubsAddr = append(sf.SubsAddr, rssURL)
	data, _ = json.Marshal(sf)
	return c.Set(sourceforgeKey, data)
}
