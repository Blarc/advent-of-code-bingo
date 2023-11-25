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

type BingoBoardDto struct {
	ID         uint           `json:"id"`
	Name       string         `json:"name"`
	BingoCards []BingoCardDto `json:"bingo_cards"`
}

func (b *BingoBoard) MapToDto() BingoBoardDto {
	var bingoCardDtos []BingoCardDto
	for _, bingoCard := range b.BingoCards {
		bingoCardDtos = append(bingoCardDtos, bingoCard.MapToDto())
	}

	return BingoBoardDto{
		ID:         b.ID,
		Name:       b.Name,
		BingoCards: bingoCardDtos,
	}
}
