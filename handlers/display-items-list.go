package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Display_items_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()
	var list models.List

	if err != nil {
		return msg, err
	}

	result := db.Where("name = ?", update.Message.Text).Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).Preload("Items").First(&list)

	if result.Error != nil {
		return msg, result.Error
	}

	err = session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.ItemsListDisplayed)

	if err != nil {
		return msg, err
	}

	var text = ""

	for index, item := range list.Items {
		text += item.Name

		if index < len(list.Items)-1 {
			text += "\n"
		}
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Main menu"),
		),
	)

	return msg, err

}
