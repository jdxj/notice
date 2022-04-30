package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

const (
	configPathEnv = "NOTICE_CONFIG"
)

var (
	cw = configWrapper{
		TelegramBot: &TelegramBot,
		Github:      &Github,
	}

	TelegramBot telegramBot
	Github      github
)

var (
	ErrConfigPathEnvNotFound = errors.New("config path env not found")
)

type configWrapper struct {
	TelegramBot *telegramBot `mapstructure:"telegram_bot"`
	Github      *github      `mapstructure:"github"`
}

type telegramBot struct {
	APIToken string `mapstructure:"api_token"`
	ChatID   int64  `mapstructure:"chat_id"`
}

type github struct {
	PersonalAccessToken string `mapstructure:"personal_access_token"`
}

func init() {
	configPath := os.Getenv(configPathEnv)
	if configPath == "" {
		panic(ErrConfigPathEnvNotFound)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cw)
	if err != nil {
		panic(err)
	}
}
