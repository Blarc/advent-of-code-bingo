package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username   string      `json:"username"`
	BingoCards []BingoCard `gorm:"many2many:user_bingo_cards;" json:"bingo_cards"`
}
