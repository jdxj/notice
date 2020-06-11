package config

import "encoding/json"

var (
	EmailKey = []byte("email")
)

type Email struct {
	Addr  string `json:"addr"`
	Token string `json:"token"`
}

func (c *Cache) GetEmail() (*Email, error) {
	data, err := c.Get(EmailKey)
	if err != nil {
		return nil, err
	}

	email := &Email{}
	return email, json.Unmarshal(data, email)
}

func (c *Cache) SetEmail(email *Email) error {
	data, _ := json.Marshal(email)
	return c.Set(EmailKey, data)
}
