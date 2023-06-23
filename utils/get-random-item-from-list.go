package utils

import (
	"math/rand"

	"github.com/dkonasov/what-to-watch-bot-serverless/models"
	"github.com/dkonasov/what-to-watch-bot-serverless/persistence"
)

func Get_random_item_from_list(list_id uint) (models.Item, error) {
	var items []models.Item
	var item models.Item
	db, err := persistence.Get_db()

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
