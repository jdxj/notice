package config

import (
	"encoding/json"
	"fmt"

	"github.com/dgraph-io/badger/v2"
)

var (
	ryfKey = []byte("ruanyifeng")
)

type RYF struct {
	Users []string `json:"users"`
}

func (ryf *RYF) String() string {
	data, _ := json.MarshalIndent(ryf, "", "    ")
	return fmt.Sprintf("%s", data)
}

func (c *Cache) GetRYF() (*RYF, error) {
	data, err := c.Get(ryfKey)
	if err != nil {
		return nil, err
	}

	ryf := &RYF{}
	return ryf, json.Unmarshal(data, ryf)
}

func (c *Cache) SetRYF(emailAddr string) error {
	if emailAddr == "" {
		return nil
	}

	ryf := &RYF{}
	ryf.Users = append(ryf.Users, emailAddr)
	data, _ := json.Marshal(ryf)
	return c.Set(ryfKey, data)
}

func (c *Cache) AddRYF(emailAddr string) error {
	ryf, err := c.GetRYF()
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("add ryf failed: %s", err)
		}
		return c.SetRYF(emailAddr)
	}

	ryf.Users = append(ryf.Users, emailAddr)
	data, _ := json.Marshal(ryf)
	return c.Set(ryfKey, data)
}
