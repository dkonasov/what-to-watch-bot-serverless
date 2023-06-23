package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Create_new_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()
	var list_id uint

	if err == nil {
		list_id, err = session.Get_active_list(strconv.FormatInt(update.Message.Chat.ID, 10))
	}

	if err == nil {
		result := db.Create(&models.Item{Name: update.Message.Text, ListID: list_id})

		if result.Error != nil {
			return msg, result.Error
		}

		return Get_main_menu(update)
	}

	return msg, err
}
