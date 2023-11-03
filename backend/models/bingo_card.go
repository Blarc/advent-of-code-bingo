package models

import "gorm.io/gorm"

type BingoCard struct {
	gorm.Model
	Description string `json:"description"`
}
