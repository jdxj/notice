package config

import (
	"encoding/json"
	"fmt"
)

var (
	emailKey = []byte("email")
)

type Email struct {
	Addr  string `json:"addr"`
	Token string `json:"token"`
}

func (c *Cache) GetEmail() (*Email, error) {
	data, err := c.Get(emailKey)
	if err != nil {
		return nil, fmt.Errorf("can not get data, key: %s, err: %s",
			emailKey, err)
	}

	email := &Email{}
	return email, json.Unmarshal(data, email)
}

func (c *Cache) SetEmail(email *Email) error {
	data, _ := json.Marshal(email)
	return c.Set(emailKey, data)
}
