package main

import (
	"strconv"

	"github.com/dkonasov/what-to-watch-bot-serverless/handlers"
	"github.com/dkonasov/what-to-watch-bot-serverless/session"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func process_message(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	step, err := session.Get_step(strconv.FormatInt(update.Message.Chat.ID, 10))
	var msg tgbotapi.MessageConfig

	if err != nil {
		return msg, err
	}

	switch step {
	case session.CreateList:
		msg, err = handlers.Create_list(update)
	case session.MainMenu:
		switch update.Message.Text {
		case "Create list":
			msg, err = handlers.Get_create_list(update)
		case "Add item to list":
			msg, err = handlers.Get_add_to_list(update)
		case "Get item from list":
			msg, err = handlers.Display_lists_for_withdraw(update)
		default:
			msg, err = handlers.Get_main_menu(update)
		}
	case session.AddItemToList:
		msg, err = handlers.Get_new_item_name(update)
	case session.GetNewItemName:
		msg, err = handlers.Create_new_item(update)
	case session.SelectListForItemWithdraw:
		msg, err = handlers.Get_random_item(update)
	case session.DisplayRandowItem:
		if update.Message.Text == "Take" {
			msg, err = handlers.Take_item(update)
		} else {
			msg, err = handlers.Get_another_random_item(update)
		}
	default:
		msg, err = handlers.Get_main_menu(update)
	}

	return msg, err
}
