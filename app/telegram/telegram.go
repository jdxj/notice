package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jdxj/notice/config"
)

var (
	bot       *tgbotapi.BotAPI
	botConfig config.TelegramBot
)

func Init() {
	var err error

	botConfig = config.GetTelegramBot()
	bot, err = tgbotapi.NewBotAPI(botConfig.APIToken)
	if err != nil {
		panic(err)
	}
}

func SendMessage(txt string) (err error) {
	msg := tgbotapi.NewMessage(botConfig.ChatID, txt)
	_, err = bot.Send(msg)
	return
}
