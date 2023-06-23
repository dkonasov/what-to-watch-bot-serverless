package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	"github.com/dkonasov/what-to-watch-bot-serverless/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Get_another_random_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	var item models.Item
	list_id, err := session.Get_active_list(strconv.FormatInt(update.Message.Chat.ID, 10))

	if err != nil {
		return msg, err
	}

	item, err = utils.Get_random_item_from_list(list_id)

	if err != nil {
		return msg, err
	}

	err = session.Set_active_item(strconv.FormatInt(update.Message.Chat.ID, 10), item.ID)

	if err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, item.Name)
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
