package handlers

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"
	"gorm.io/gorm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Get_main_menu(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	err := session.Set_step(strconv.FormatInt(update.Message.Chat.ID, 10), session.MainMenu)
	has_lists := false
	has_items := false
	var lists []models.List
	var db *gorm.DB

	if err == nil {
		db, err = persistence.Get_db()

		if err == nil {
			query_result := db.Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).Preload("Items").Find(&lists)
			err = query_result.Error
			has_lists = query_result.RowsAffected > 0
		}
	}

	if err == nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Hello! This is what to watch bot, a bot, that helps you to maintain your backlog lists for movies/TV shows/games etc.")
		keyboard := tgbotapi.NewReplyKeyboard(
			tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Create list"),
			),
		)

		keyboard.OneTimeKeyboard = true

		if has_lists {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton("Add item to list"),
			))

			for _, list := range lists {
				if len(list.Items) > 0 {
					has_items = true
					break
				}
			}

			if has_items {
				keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Get item from list"),
				), tgbotapi.NewKeyboardButtonRow(
					tgbotapi.NewKeyboardButton("Show items"),
				))
			}

		}

		msg.ReplyMarkup = keyboard
	}

	return msg, err
}
