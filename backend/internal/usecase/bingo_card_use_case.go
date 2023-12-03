package usecase

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/google/uuid"
)

type BingoCardUseCase struct {
	repo BingoCardRepo
}

func NewBingoCardUseCase(repo BingoCardRepo) *BingoCardUseCase {
	return &BingoCardUseCase{repo: repo}
}

func (u *BingoCardUseCase) GetBingoCards() ([]entity.BingoCardDto, error) {
	bingoCards, err := u.repo.GetBingoCards()
	if err != nil {
		return nil, fmt.Errorf("BingoCardUseCase - GetBingoCards - u.bingoBoardRepo.GetBingoCards: %w", err)
	}

	var bingoCardDtos []entity.BingoCardDto
	for _, bingoCard := range bingoCards {
		bingoCardDtos = append(bingoCardDtos, bingoCard.MapToDto())
	}

	return bingoCardDtos, nil
}

func (u *BingoCardUseCase) ClickBingoCard(user *entity.User, id uuid.UUID) error {
	bingoCard, err := u.repo.GetBingoCard(id)
	if err != nil {
		return fmt.Errorf("BingoCardUseCase - ClickBingoCard - u.bingoBoardRepo.GetBingoCard: %w", err)
	}

	if user.HasBingoCard(&bingoCard) {
		err := u.repo.RemoveBingoCardFromUser(user, &bingoCard)
		if err != nil {
			return fmt.Errorf("BingoCardUseCase - ClickBingoCard - u.bingoBoardRepo.RemoveBingoCardFromUser: %w", err)
		}
	} else {
		err := u.repo.AddBingoCardToUser(user, &bingoCard)
		if err != nil {
			return fmt.Errorf("BingoCardUseCase - ClickBingoCard - u.bingoBoardRepo.AddBingoCardToUser: %w", err)
		}
	}
	return nil
}
