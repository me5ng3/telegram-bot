package main

import (
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config := LoadConfig()
	client := http.Client{Timeout: time.Second * time.Duration(rand.Int31n(int32(config.Timeout)))}
	bot, err := telegram.NewBotAPIWithClient(config.Token, &client)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = config.Debug

	cmdHandler := newCommandHandler(bot, 10, config)
	cmdHandler.RegisterCommand("help", false, Help)

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = int(time.Second * 5)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := telegram.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		if update.Message.IsCommand() {
			go cmdHandler.Check(update.Message.Command(), strings.Split(update.Message.CommandArguments(), " "), &update)
		} else {
			bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Please use /help for a list of all the available commands"))
		}
	}
}
