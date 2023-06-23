package models

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	Name   string
	ListID uint
	List   List
}
