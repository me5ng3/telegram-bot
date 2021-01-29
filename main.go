package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	cron "github.com/robfig/cron"
)

func main() {
	scheduler := cron.New()
	defer scheduler.Stop()

	client := http.Client{Timeout: time.Second * 5}
	bot, err := telegram.NewBotAPIWithClient("1615998279:AAG8QzNbtO61mtnF5AKTz6qivnuBWzjasPY", &client)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	scheduler.AddFunc("0 1 * * * *", func() { fmt.Println("Every hour on the half hour") })

	u := telegram.NewUpdate(0)
	u.Timeout = int(time.Second * 5)

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := telegram.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		if update.Message.IsCommand() {
			msg := telegram.NewMessage(update.Message.Chat.ID, "")
			switch update.Message.Command() {
			case "help":
				msg.Text = "type /sayhi or /status."
			case "sayhi":
				msg.Text = "Hi :)"
			case "status":
				msg.Text = "I'm ok."
			case "withArgument":
				msg.Text = "You supplied the following argument: " + update.Message.CommandArguments()
			case "html":
				msg.ParseMode = "html"
				msg.Text = "This will be interpreted as HTML, click <a href=\"https://www.example.com\">here</a>"
			default:
				msg.Text = "This is not a valid command"
			}
			bot.Send(msg)
		} else {
			bot.Send(telegram.NewMessage(update.Message.Chat.ID, "Please use /help for a list of all the available commands"))
		}
	}
}
