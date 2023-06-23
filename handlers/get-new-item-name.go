package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Get_new_item_name(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()
	var list models.List

	if err == nil {
		result := db.Where("name = ?", update.Message.Text).Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).First(&list)
		err = result.Error
	}

	if err == nil {
		err = session.Set_active_list(strconv.FormatInt(update.Message.Chat.ID, 10), list.ID)
	}

	if err == nil {
		err = session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.GetNewItemName)
	}

	if err == nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Enter item name")
	}

	return msg, err
}
