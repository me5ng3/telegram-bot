package main

import (
	"fmt"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CommandHandler struct {
	bot      *telegram.BotAPI
	commands map[string]*Command
	logger   <-chan string
}

type Command struct {
	name           string
	onlyRegistered bool
	function       func(*CommandHandler, *telegram.Update, []string)
}

func (cmdHandler *CommandHandler) startLogging(chatSize int) {
	logger := make(chan string, chatSize)
	cmdHandler.logger = logger

	go func(logger <-chan string) {
		for {
			for message := range logger {
				fmt.Println(message) // LOG FORMATTING! <DATE:HOUR> MESSAGE
				// WRITE LOGS TO POSTGRESQL
			}
		}
	}(logger)
}

func newCommandHandler(bot *telegram.BotAPI) *CommandHandler {
	return &CommandHandler{bot: bot, commands: make(map[string]*Command)}
}

func (cmdHandler *CommandHandler) RegisterCommand(name string, onlyRegistered bool, function func(*CommandHandler, *telegram.Update, []string)) {
	cmdHandler.commands[name] = &Command{name, onlyRegistered, function}
}

func (cmdHandler *CommandHandler) Check(commandName string, args []string, u *telegram.Update) {
	if command, ok := cmdHandler.commands[commandName]; ok {
		// If message author in database: continue, else: access forbidden.
		command.function(cmdHandler, u, args)
	} else {
		cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "Not a valid command. Please use /help for a list of all the commands."))
	}
}
