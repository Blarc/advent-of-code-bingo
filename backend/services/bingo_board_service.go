package services

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/Blarc/advent-of-code-bingo/repo"
)

type BingoBoardService interface {
	FindBingoBoard(shortUuid string) (*models.BingoBoard, error)
}

type bingoBoardService struct {
	bingoBoardRepo repo.BingoBoardRepo
}

func NewBingoBoardService(bingoBoardRepo repo.BingoBoardRepo) BingoBoardService {
	return &bingoBoardService{bingoBoardRepo: bingoBoardRepo}
}

func (s *bingoBoardService) FindBingoBoard(shortUuid string) (*models.BingoBoard, error) {
	// Find the bingo board by the first 16 characters of the ID
	bingoBoard, err := s.bingoBoardRepo.FindBingoBoard(shortUuid)
	if err != nil {
		return nil, err
	}
	return bingoBoard, nil
}
