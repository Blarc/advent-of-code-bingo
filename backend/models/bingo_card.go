package models

import "gorm.io/gorm"

type BingoCard struct {
	gorm.Model
	ID          uint   `gorm:"primarykey" json:"id"`
	Description string `json:"description"`
}
