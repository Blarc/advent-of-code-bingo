package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BingoCard struct {
	gorm.Model
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Description string    `json:"description"`
	Public      bool      `json:"public"`
}

type BingoCardDto struct {
	ID          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	UserCount   uint      `json:"user_count"`
	Selected    bool      `json:"selected"`
}

func (b *BingoCard) MapToDto() BingoCardDto {
	return BingoCardDto{
		ID:          b.ID,
		Description: b.Description,
	}
}

type BingoCardId struct {
	ID string `uri:"id" binding:"required"`
}
