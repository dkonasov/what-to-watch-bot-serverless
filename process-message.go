package main

import (
	"strconv"

	"gorm.io/gorm"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func process_message(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	step, err := get_step(strconv.FormatInt(update.Message.Chat.ID, 10))
	var msg tgbotapi.MessageConfig

	if err != nil {
		return msg, err
	}

	switch step {
	case CreateList:
		msg, err = create_list(update)
	case MainMenu:
		switch update.Message.Text {
		case "Create list":
			msg, err = get_create_list(update)
		case "Add item to list":
			msg, err = get_add_to_list(update)
		case "Get item from list":
			msg, err = display_lists_for_withdraw(update)
		default:
			msg, err = get_main_menu(update)
		}
	case AddItemToList:
		msg, err = get_new_item_name(update)
	case GetNewItemName:
		msg, err = create_new_item(update)
	case SelectListForItemWithdraw:
		msg, err = get_random_item(update)
	case DisplayRandowItem:
		if update.Message.Text == "Take" {
			msg, err = take_item(update)
		} else {
			msg, err = get_another_random_item(update)
		}
	default:
		msg, err = get_main_menu(update)
	}

	return msg, err
}

func get_main_menu(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	err := set_step(strconv.FormatInt(update.Message.Chat.ID, 10), MainMenu)
	has_lists := false
	has_items := false
	var lists []List
	var db *gorm.DB

	if err == nil {
		db, err = get_db()

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
				))
			}

		}

		msg.ReplyMarkup = keyboard
	}

	return msg, err
}

func get_create_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig

	err := set_step(strconv.FormatInt(update.Message.Chat.ID, 10), CreateList)

	if err == nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Enter name of the list")
	}

	return msg, err
}

func create_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := get_db()

	if err != nil {
		return msg, err
	}

	db.Create(&List{User: strconv.FormatInt(update.Message.Chat.ID, 10), Name: update.Message.Text})

	return get_main_menu(update)
}

func get_add_to_list(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var lists []List
	var msg tgbotapi.MessageConfig
	db, err := get_db()

	if err == nil {
		result := db.Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).Find(&lists)
		err = result.Error
	}

	if err == nil {
		err = set_step(strconv.FormatInt(update.Message.Chat.ID, 10), AddItemToList)
	}

	if err == nil {
		keyboard := tgbotapi.NewReplyKeyboard()
		keyboard.OneTimeKeyboard = true

		for _, list := range lists {
			keyboard.Keyboard = append(keyboard.Keyboard, tgbotapi.NewKeyboardButtonRow(
				tgbotapi.NewKeyboardButton(list.Name),
			))
		}

		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Select list")
		msg.ReplyMarkup = keyboard
	}

	return msg, err
}

func get_new_item_name(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := get_db()
	var list List

	if err == nil {
		result := db.Where("name = ?", update.Message.Text).Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).First(&list)
		err = result.Error
	}

	if err == nil {
		err = set_active_list(strconv.FormatInt(update.Message.Chat.ID, 10), list.ID)
	}

	if err == nil {
		err = set_step(strconv.FormatInt(update.Message.Chat.ID, 10), GetNewItemName)
	}

	if err == nil {
		msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Enter item name")
	}

	return msg, err
}

func create_new_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := get_db()
	var list_id uint

	if err == nil {
		list_id, err = get_active_list(strconv.FormatInt(update.Message.Chat.ID, 10))
	}

	if err == nil {
		result := db.Create(&Item{Name: update.Message.Text, ListID: list_id})

		if result.Error != nil {
			return msg, result.Error
		}

		return get_main_menu(update)
	}

	return msg, err
}

func display_lists_for_withdraw(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var lists []List
	var msg tgbotapi.MessageConfig
	db, err := get_db()

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

	err = set_step(strconv.FormatInt(update.Message.Chat.ID, 10), SelectListForItemWithdraw)

	if err != nil {
		return msg, err
	}

	msg = tgbotapi.NewMessage(update.Message.Chat.ID, "Select list")
	msg.ReplyMarkup = keyboard

	return msg, err
}

func get_random_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	db, err := get_db()
	var list List

	if err != nil {
		return msg, err
	}

	result := db.Where("name = ?", update.Message.Text).Where("\"user\" = ?", strconv.FormatInt(update.Message.Chat.ID, 10)).First(&list)

	if result.Error != nil {
		return msg, result.Error
	}

	var selected_item Item
	selected_item, err = get_random_item_from_list(list.ID)

	if err != nil {
		return msg, err
	}

	err = set_active_list(strconv.FormatInt(update.Message.Chat.ID, 10), list.ID)

	if err != nil {
		return msg, err
	}

	err = set_active_item(strconv.FormatInt(update.Message.Chat.ID, 10), selected_item.ID)

	if err != nil {
		return msg, err
	}

	err = set_step(strconv.FormatInt(update.Message.Chat.ID, 10), DisplayRandowItem)

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

func get_another_random_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	var item Item
	list_id, err := get_active_list(strconv.FormatInt(update.Message.Chat.ID, 10))

	if err != nil {
		return msg, err
	}

	item, err = get_random_item_from_list(list_id)

	if err != nil {
		return msg, err
	}

	err = set_active_item(strconv.FormatInt(update.Message.Chat.ID, 10), item.ID)

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

func take_item(update *tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	var msg tgbotapi.MessageConfig
	var item Item
	id, err := get_active_item(strconv.FormatInt(update.Message.Chat.ID, 10))

	if err != nil {
		return msg, err
	}

	db, err := get_db()

	if err != nil {
		return msg, err
	}

	result := db.Where("id = ?", id).Delete(&item)

	if result.Error != nil {
		return msg, result.Error
	}

	return get_main_menu(update)
}
