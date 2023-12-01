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
	result := b.db.
		Preload("Users").
		First(&bingoBoard, "substring(id::text, 1, 16) = ?", shortUuid)
	if result.Error != nil {
		return nil, result.Error
	}
	return &bingoBoard, nil
}

// GetBingoCardsWithCount returns all bingo cards for a board and how many users have them
func (b *bingoBoardRepo) GetBingoCardsWithCount(shortUuid string) ([]models.BingoCardDto, error) {
	var bingoCards []models.BingoCardDto
	result := models.DB.Table("bingo_board_bingo_card").
		Select("bingo_cards.id, bingo_cards.description, count(distinct user_id) as user_count").
		Joins("left join bingo_cards ON bingo_cards.id = bingo_board_bingo_card.bingo_card_id").
		Joins(`left join user_bingo_card ON user_bingo_card.bingo_card_id = bingo_board_bingo_card.bingo_card_id AND user_bingo_card.user_id IN (
						SELECT user_id
						FROM user_bingo_board
						WHERE substring(bingo_board_id::text, 1, 16) = ?
					)`, shortUuid).
		Where("substring(bingo_board_bingo_card.bingo_board_id::text, 1, 16) = ?", shortUuid).
		Group("bingo_cards.id").
		Scan(&bingoCards)

	if result.Error != nil {
		return nil, result.Error
	}

	return bingoCards, nil
}
