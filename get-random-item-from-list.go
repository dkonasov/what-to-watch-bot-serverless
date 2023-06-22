package main

import "math/rand"

func get_random_item_from_list(list_id uint) (Item, error) {
	var items []Item
	var item Item
	db, err := get_db()

	if err != nil {
		return item, err
	}

	result := db.Where("list_id = ?", list_id).Find(&items)

	if result.Error != nil {
		return item, result.Error
	}

	item = items[rand.Intn(int(result.RowsAffected))]

	return item, err
}
