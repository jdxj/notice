package main

import (
	"log"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/jdxj/notice/config"
)

var (
	bot *tg.BotAPI
)

func init() {
	var err error
	bot, err = tg.NewBotAPI(config.TelegramBot.APIToken)
	if err != nil {
		panic(err)
	}
	log.Printf("init telegram ok")
}

func SendMessage(txt string) (err error) {
	msg := tg.NewMessage(config.TelegramBot.ChatID, txt)
	_, err = bot.Send(msg)
	return
}
