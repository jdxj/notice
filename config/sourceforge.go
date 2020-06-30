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

func (sf *Sourceforge) String() string {
	data, _ := json.MarshalIndent(sf, "", "    ")
	return fmt.Sprintf("%s", data)
}

func GetSourceforge() (*Sourceforge, error) {
	data, err := get(sourceforgeKey)
	if err != nil {
		return nil, err
	}
	sf := &Sourceforge{}
	return sf, json.Unmarshal(data, sf)
}

func SetSFSubsAddr(rssURL string) error {
	sf := &Sourceforge{}
	sf.SubsAddr = append(sf.SubsAddr, rssURL)
	data, _ := json.Marshal(sf)
	return set(sourceforgeKey, data)
}

func AddSFSubsAddr(rssURL string) error {
	// url 检查
	if rssURL == "" {
		return fmt.Errorf("rss address invalid")
	}

	sf, err := GetSourceforge()
	if err != nil {
		if err != badger.ErrKeyNotFound {
			return fmt.Errorf("add subscription address failed: %s", err)
		}
		return SetSFSubsAddr(rssURL)
	}

	// 查重
	for _, v := range sf.SubsAddr {
		if v == rssURL {
			return fmt.Errorf("duplicate subscription address: %s", rssURL)
		}
	}
	sf.SubsAddr = append(sf.SubsAddr, rssURL)
	data, _ := json.Marshal(sf)
	return set(sourceforgeKey, data)
}
