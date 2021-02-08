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

	// Load configuration file (config.json)
	conf := loadConfig()

	// start telegram bot with http client
	client := http.Client{Timeout: time.Second * time.Duration(rand.Int31n(int32(conf.Timeout)))}
	bot, err := telegram.NewBotAPIWithClient(conf.Token, &client)
	if err != nil {
		log.Panic(err)
	}

	// set debug variable
	bot.Debug = conf.Debug

	// initialize command handler, pass bot and configuration file
	// initialize commands
	cmdHandler := NewCommandHandler(bot, conf)
	cmdHandler.commands["help"] = &Command{"help", false, Help}
	cmdHandler.commands["corona"] = &Command{"corona", false, CoronaUpdate}

	// online
	log.Printf("Authorized on account %s", bot.Self.UserName)

	// set update timer
	u := telegram.NewUpdate(0)
	u.Timeout = int(time.Second * 5)

	// get updates
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	// loop through latest updates. If update is command then pass it through a go routine to
	// the command handler along with command arguments and message pointer.
	for update := range updates {
		msg := telegram.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID
		if update.Message.IsCommand() {
			go cmdHandler.Check(update.Message.Command(), strings.Split(update.Message.CommandArguments(), " "), &update)
		} else {
			bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Please use /help for a list of all the available commands"))
		}
	}
}
