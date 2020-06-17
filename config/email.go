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

func (e *Email) String() string {
	data, _ := json.MarshalIndent(e, "", "    ")
	return fmt.Sprintf("%s", data)
}

func (c *Cache) GetEmail() (*Email, error) {
	data, err := c.Get(emailKey)
	if err != nil {
		return nil, err
	}

	email := &Email{}
	return email, json.Unmarshal(data, email)
}

func (c *Cache) SetEmail(email *Email) error {
	data, _ := json.Marshal(email)
	return c.Set(emailKey, data)
}
