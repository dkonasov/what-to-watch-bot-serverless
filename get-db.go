package main

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func get_db() (*gorm.DB, error) {
	var (
		host     = os.Getenv("PG_HOST")
		port     = os.Getenv("PG_PORT")
		user     = os.Getenv("PG_USER")
		password = os.Getenv("PG_PASSWORD")
		dbname   = os.Getenv("PG_DB")
	)

	connstring := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s target_session_attrs=read-write",
		host, port, dbname, user, password)
	db, err := gorm.Open(postgres.Open(connstring), &gorm.Config{})

	db.AutoMigrate(&List{})
	db.AutoMigrate(&Item{})

	return db, err
}
