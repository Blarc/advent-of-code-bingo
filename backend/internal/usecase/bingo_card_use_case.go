package usecase

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
)

type BingoCardUseCase struct {
	repo BingoCardRepo
}

func New(repo BingoCardRepo) BingoCardUseCase {
	return BingoCardUseCase{repo: repo}
}

func (u *BingoCardUseCase) GetBingoCards() ([]entity.BingoCardDto, error) {
	bingoCardsDto, err := u.repo.GetBingoCards()
	if err != nil {
		return nil, fmt.Errorf("BingoCardUseCase - GetBingoCards - u.bingoBoardRepo.GetBingoCards: %w", err)
	}
	return bingoCardsDto, nil
}
