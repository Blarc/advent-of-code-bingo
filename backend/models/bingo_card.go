package models

import "gorm.io/gorm"

type BingoCard struct {
	gorm.Model
	Description string `json:"description"`
}

type BingoCardDto struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	UserCount   uint   `json:"user_count"`
	Selected    bool   `json:"selected"`
}

func (b *BingoCard) MapToDto() BingoCardDto {
	return BingoCardDto{
		ID:          b.ID,
		Description: b.Description,
	}
}
