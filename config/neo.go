package config

import (
	"encoding/json"
	"fmt"
)

var (
	neoKey = []byte("neo")
)

type Neo struct {
	Host     string `json:"host"`
	Domain   string `json:"domain"`
	Cookies  string `json:"cookies"`
	Services string `json:"services"`
	User     string `json:"user"`
}

func (c *Cache) GetNeo() (*Neo, error) {
	data, err := c.Get(neoKey)
	if err != nil {
		return nil, fmt.Errorf("can not get data, key: %s, err: %s",
			neoKey, err)
	}

	neo := &Neo{}
	return neo, json.Unmarshal(data, neo)
}

func (c *Cache) SetNeo(neo *Neo) error {
	data, _ := json.Marshal(neo)
	return c.Set(neoKey, data)
}
