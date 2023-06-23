package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Create_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := persistence.Get_db()

	if err != nil {
		return msg, err
	}

	db.Create(&models.List{User: strconv.FormatInt(update.Message.Chat.ID, 10), Name: update.Message.Text})

	return Get_main_menu(update)
}
