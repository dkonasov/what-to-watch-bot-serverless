package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Get_create_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig

	err := session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.CreateList)

	if err == nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Enter name of the list")
	}

	return msg, err
}
