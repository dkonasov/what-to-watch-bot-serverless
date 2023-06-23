package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	"github.com/dkonasov/what-to-watch-bot-serverless/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Get_random_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()
	var list models.List

	if err != nil {
		return msg, err
	}

	result := db.Where("name = ?", update.Message.Text).Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).First(&list)

	if result.Error != nil {
		return msg, result.Error
	}

	var selected_item models.Item
	selected_item, err = utils.Get_random_item_from_list(list.ID)

	if err != nil {
		return msg, err
	}

	err = session.Set_active_list(strconv.FormatInt(update.Message.Chat.ID, 10), list.ID)

	if err != nil {
		return msg, err
	}

	err = session.Set_active_item(strconv.FormatInt(update.Message.Chat.ID, 10), selected_item.ID)

	if err != nil {
		return msg, err
	}

	err = session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.DisplayRandowItem)

	if err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, selected_item.Name)
	msg.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Take"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Get another"),
		),
	)

	return msg, err

}
