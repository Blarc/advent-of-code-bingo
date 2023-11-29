package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BingoBoard struct {
	gorm.Model
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid()"`
	Name       string      `json:"name"`
	OwnerId    uuid.UUID   `gorm:"type:uuid;unique" json:"owner_id"`
	Users      []User      `gorm:"many2many:user_bingo_board;" json:"users"`
	BingoCards []BingoCard `gorm:"many2many:bingo_board_bingo_card;" json:"bingo_cards"`
}

type BingoBoardDto struct {
	ShortUUID  string         `gorm:"size:16" json:"short_uuid"`
	Name       string         `json:"name"`
	BingoCards []BingoCardDto `json:"bingo_cards"`
}

func (b *BingoBoard) MapToDto() BingoBoardDto {
	var bingoCardDtos []BingoCardDto
	for _, bingoCard := range b.BingoCards {
		bingoCardDtos = append(bingoCardDtos, bingoCard.MapToDto())
	}

	return BingoBoardDto{
		ShortUUID:  b.ID.String()[:16],
		Name:       b.Name,
		BingoCards: bingoCardDtos,
	}
}

type BingoBoardId struct {
	ID string `uri:"id" binding:"required"`
}
