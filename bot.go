package main

import (
	"fmt"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	token string
}

func (bot *Bot) HandleRequest(update *tgbotapi.Update) error {
	api, err := tgbotapi.NewBotAPI(bot.token)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while initing API: %s", err.Error())

		return err
	}

	msg, err := process_message(update)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while processing message: %s", err.Error())

		return err
	}

	_, err = api.Send(msg)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error while sending message: %s", err.Error())

		return err
	}

	return nil
}
