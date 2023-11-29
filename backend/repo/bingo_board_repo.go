package repo

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"gorm.io/gorm"
)

type BingoBoardRepo interface {
	FindBingoBoard(id string) (*models.BingoBoard, error)
	GetBingoCardsWithCount(uuid string) ([]models.BingoCardDto, error)
}

type bingoBoardRepo struct {
	db *gorm.DB
}

func NewBingoBoardRepo(db *gorm.DB) BingoBoardRepo {
	return &bingoBoardRepo{db: db}
}

func (b *bingoBoardRepo) FindBingoBoard(shortUuid string) (*models.BingoBoard, error) {
	var bingoBoard models.BingoBoard
	result := b.db.First(&bingoBoard, "substring(id::text, 1, 16) = ?", shortUuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bingoBoard, nil
}

// GetBingoCardsWithCount returns all bingo cards for a board and how many users have them
func (b *bingoBoardRepo) GetBingoCardsWithCount(shortUuid string) ([]models.BingoCardDto, error) {
	var bingoCards []models.BingoCardDto
	result := models.DB.Table("bingo_cards").
		Select("bingo_cards.id, bingo_cards.description, count(user_bingo_card.user_id) as user_count").
		Joins("left join bingo_board_bingo_card ON bingo_board_bingo_card.bingo_card_id = bingo_cards.id").
		Joins("left join user_bingo_card ON user_bingo_card.bingo_card_id = bingo_cards.id").
		Group("bingo_board_bingo_card.bingo_board_id, bingo_cards.id").
		Find(&bingoCards, "substring(bingo_board_bingo_card.bingo_board_id::text, 1, 16) = ?", shortUuid)

	if result.Error != nil {
		return nil, result.Error
	}

	return bingoCards, nil
}
