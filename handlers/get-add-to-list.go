package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Get_add_to_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var lists []models.List
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()

	if err == nil {
		result := db.Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).Find(&lists)
		err = result.Error
	}

	if err == nil {
		err = session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.AddItemToList)
	}

	if err == nil {
		keyboard := tgbotapi.NewReplyKeyboard()
		keyboard.OneTimeKeyboard = true

		for _, list := range lists {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(list.Name),
			))
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Select list")
		msg.ReplyMarkup = keyboard
	}

	return msg, err
}
