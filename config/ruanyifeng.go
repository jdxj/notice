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

func GetRYF() (*RYF, error) {
	data, err := get(ryfKey)
	if err != nil {
		return nil, err
	}
	ryf := &RYF{}
	return ryf, json.Unmarshal(data, ryf)
}

func SetRYF(emailAddr string) error {
	if emailAddr == "" {
		return nil
	}

	ryf := &RYF{}
	ryf.Users = append(ryf.Users, emailAddr)
	data, _ := json.Marshal(ryf)
	return set(ryfKey, data)
}

func AddRYF(emailAddr string) error {
	ryf, err := GetRYF()
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("add ryf failed: %s", err)
		}
		return SetRYF(emailAddr)
	}

	ryf.Users = append(ryf.Users, emailAddr)
	data, _ := json.Marshal(ryf)
	return set(ryfKey, data)
}
