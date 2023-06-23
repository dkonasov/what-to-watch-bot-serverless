package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Take_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	var item models.Item
	id, err := session.Get_active_item(strconv.FormatInt(update.Message.Chat.ID, 10))

	if err != nil {
		return msg, err
	}

	db, err := persistence.Get_db()

	if err != nil {
		return msg, err
	}

	result := db.Where("id = ?", id).Delete(&item)

	if result.Error != nil {
		return msg, result.Error
	}

	return Get_main_menu(update)
}
