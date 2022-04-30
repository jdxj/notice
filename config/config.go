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
}

func Init(paths ...string) error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	for _, p := range paths {
		viper.AddConfigPath(p)
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return viper.Unmarshal(&defaultConfig)
}

