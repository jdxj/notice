package module

import (
	"encoding/json"
	"os"
)

type Config struct {
	NeoProxy *NeoProxyConfig `json:"neoProxy"`
	Email    *EmailConfig    `json:"email"`
}

type NeoProxyConfig struct {
	URL     string `json:"url"`
	Domain  string `json:"domain"`
	Cookie  string `json:"cookie"`
	Service string `json:"service"`
}

func ReadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	return config, decoder.Decode(config)
}

type EmailConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}
