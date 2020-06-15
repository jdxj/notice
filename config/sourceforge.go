package config

import (
	"encoding/json"
	"fmt"
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

func (c *Cache) AddSubsAddr(rssURL string) error {
	if rssURL == "" {
		return fmt.Errorf("rss address invalid")
	}

	sf, err := c.GetSourceforge()
	if err != nil {
		return fmt.Errorf("add subscription address failed: %s", err)
	}

	for _, v := range sf.SubsAddr {
		if v == rssURL {
			return fmt.Errorf("duplicate subscription address: %s", rssURL)
		}
	}

	sf.SubsAddr = append(sf.SubsAddr, rssURL)
	data, _ := json.Marshal(sf)
	return c.Set(sourceforgeKey, data)
}
