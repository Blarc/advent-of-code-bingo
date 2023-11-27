package repo

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"gorm.io/gorm"
)

type BingoBoardRepo interface {
	FindBingoBoard(id string) (*models.BingoBoard, error)
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
