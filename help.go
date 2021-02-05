package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func Help(cmdHandler *CommandHandler, u *telegram.Update, args []string) {
	cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "To be updated soon."))
}
