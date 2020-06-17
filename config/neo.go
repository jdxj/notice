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

func (neo *Neo) String() string {
	data, _ := json.MarshalIndent(neo, "", "    ")
	return fmt.Sprintf("%s", data)
}

func (c *Cache) GetNeo() (*Neo, error) {
	data, err := c.Get(neoKey)
	if err != nil {
		return nil, err
	}

	neo := &Neo{}
	return neo, json.Unmarshal(data, neo)
}

func (c *Cache) SetNeo(neo *Neo) error {
	data, _ := json.Marshal(neo)
	return c.Set(neoKey, data)
}
