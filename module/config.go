package module

import (
	"encoding/json"
	"os"
)

const (
	configPath = "config.json"
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
	Email   string `json:"email"`
}

type EmailConfig struct {
	Address string `json:"address"`
	Token   string `json:"token"`
}

func ReadConfig() (*Config, error) {
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	return config, decoder.Decode(config)
}

func WriteConfig(config *Config) error {
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return err
	}
	defer file.Close()
	defer file.Sync()

	_, err = file.Write(data)
	return err
}
