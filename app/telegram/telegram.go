package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jdxj/notice/config"
)

var (
	bot *tgbotapi.BotAPI
)

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(config.TelegramBot.APIToken)
	if err != nil {
		panic(err)
	}
}

func SendMessage(txt string) (err error) {
	msg := tgbotapi.NewMessage(config.TelegramBot.ChatID, txt)
	_, err = bot.Send(msg)
	return
}
