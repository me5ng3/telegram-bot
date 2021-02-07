package main

import (
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

type CommandHandler struct {
	bot      *telegram.BotAPI
	commands map[string]*Command
	config   *config
}

type Command struct {
	name           string
	onlyRegistered bool
	function       func(*CommandHandler, *telegram.Update, []string)
}

func newCommandHandler(bot *telegram.BotAPI, config *config) *CommandHandler {
	return &CommandHandler{bot: bot, commands: make(map[string]*Command), config: config}
}

func (cmdHandler *CommandHandler) Check(commandName string, args []string, u *telegram.Update) {
	if command, ok := cmdHandler.commands[commandName]; ok {
		// If message author in database: continue, else: access forbidden.
		command.function(cmdHandler, u, args)
	} else {
		cmdHandler.bot.Send(telegram.NewMessage(u.Message.Chat.ID, "Not a valid command. Please use /help for a list of all the commands."))
	}
}
