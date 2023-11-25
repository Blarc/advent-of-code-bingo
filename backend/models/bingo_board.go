package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BingoBoard struct {
	gorm.Model
	Name       string      `json:"name"`
	UserID     uuid.UUID   `json:"user_id"`
	BingoCards []BingoCard `gorm:"many2many:bingo_board_bingo_card;" json:"bingo_cards"`
}
