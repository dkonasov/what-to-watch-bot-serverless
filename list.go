package main

import (
	"gorm.io/gorm"
)

type List struct {
	gorm.Model
	User  string
	Name  string
	Items []Item
}
