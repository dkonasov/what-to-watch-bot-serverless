package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Display_lists_for_withdraw(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var lists []models.List
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()

	if err != nil {
		return msg, err
	}

	result := db.Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).Preload("Items").Find(&lists)

	if result.Error != nil {
		return msg, result.Error
	}

	keyboard := tgbotapi.NewReplyKeyboard()
	keyboard.OneTimeKeyboard = true

	for _, list := range lists {
		if len(list.Items) > 0 {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(list.Name),
			))
		}
	}

	err = session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.SelectListForItemWithdraw)

	if err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Select list")
	msg.ReplyMarkup = keyboard

	return msg, err
}
