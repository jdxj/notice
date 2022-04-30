package config

import "github.com/spf13/viper"

var (
	defaultConfig config
)

type config struct {
	TelegramBot TelegramBot `mapstructure:"telegram_bot"`
}

type TelegramBot struct {
	APIToken string `mapstructure:"api_token"`
	ChatID   int64  `mapstructure:"chat_id"`
}

func Init(paths ...string) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	for _, p := range paths {
		viper.AddConfigPath(p)
	}

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&defaultConfig)
	if err != nil {
		panic(err)
	}
}

func GetTelegramBot() TelegramBot {
	return defaultConfig.TelegramBot
}
