package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {

	client := http.Client{Timeout: time.Second * 5}
	bot, err := telegram.NewBotAPIWithClient("1615998279:AAG8QzNbtO61mtnF5AKTz6qivnuBWzjasPY", &client)
	if err != nil {
		log.Panic(err)
	}

	cmdHandler := newCommandHandler(bot)
	cmdHandler.startLogging(10)
	cmdHandler.RegisterCommand("help", false, Help)

	// bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = int(time.Second * 5)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := telegram.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			go cmdHandler.Check(update.Message.Command(), strings.Split(update.Message.CommandArguments(), " "), &update)
		} else {
			bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Please use /help for a list of all the available commands"))
		}
	}
}
