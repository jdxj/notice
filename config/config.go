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
		DB:          &DB,
		RSS:         &RSS,
	}

	TelegramBot telegramBot
	Github      github
	DB          db
	RSS         rss
)

var (
	ErrConfigPathEnvNotFound = errors.New("config path env not found")
)

type configWrapper struct {
	TelegramBot *telegramBot `mapstructure:"telegram_bot"`
	Github      *github      `mapstructure:"github"`
	DB          *db          `mapstructure:"db"`
	RSS         *rss         `mapstructure:"rss"`
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

type telegramBot struct {
	APIToken string `mapstructure:"api_token"`
	ChatID   int64  `mapstructure:"chat_id"`
}

type github struct {
	// https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token
	PersonalAccessToken string `mapstructure:"personal_access_token"`

	Spec string `mapstructure:"spec"`
}

type db struct {
	User string `mapstructure:"user"`
	Pass string `mapstructure:"pass"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Name string `mapstructure:"name"`
}

type rss struct {
	Spec string `mapstructure:"spec"`
}
