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
		return nil, err
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

	sf, err := c.GetSourceforge()
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("add subscription address failed: %s", err)
		}

		return c.SetSubsAddr(rssURL)
	}

	// 查重
	for _, v := range sf.SubsAddr {
		if v == rssURL {
			return fmt.Errorf("duplicate subscription address: %s", rssURL)
		}
	}

	sf.SubsAddr = append(sf.SubsAddr, rssURL)
	data, _ := json.Marshal(sf)
	return c.Set(sourceforgeKey, data)
}
