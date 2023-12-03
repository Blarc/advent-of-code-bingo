package usecase

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
)

type BingoBoardUseCase struct {
	bingoBoardRepo BingoBoardRepo
	bingoCardRepo  BingoCardRepo
}

func NewBingoBoardUseCase(bingoBoardRepo BingoBoardRepo, bingoCardRepo BingoCardRepo) *BingoBoardUseCase {
	return &BingoBoardUseCase{bingoBoardRepo, bingoCardRepo}
}

func (u *BingoBoardUseCase) GetBingoBoard(shortUuid string) (*entity.BingoBoardDto, error) {
	bingoBoard, err := u.bingoBoardRepo.GetBingoBoard(shortUuid)
	if err != nil {
		return nil, fmt.Errorf("BingoBoardUseCase - GetBingoBoard - u.bingoBoardRepo.GetBingoBoard: %w", err)
	}

	bingoCards, err := u.bingoBoardRepo.GetBingoCardsWithCount(shortUuid)
	if err != nil {
		return nil, fmt.Errorf("BingoBoardUseCase - GetBingoBoard - u.bingoBoardRepo.GetBingoCardsWithCount: %w", err)
	}

	// TODO: MapToDto returns pointer, does this work?
	bingoBoardDto := bingoBoard.MapToDto()
	bingoBoardDto.BingoCards = bingoCards
	return &bingoBoardDto, nil
}

func (u *BingoBoardUseCase) CreatePersonalBingoBoard(user *entity.User) error {
	bingoCards, err := u.bingoCardRepo.GetBingoCards()
	if err != nil {
		return fmt.Errorf("BingoBoardUseCase - CreatePersonalBingoBoard - u.bingoCardRepo.GetBingoCards: %w", err)
	}

	bingoBoard := &entity.BingoBoard{
		Name:       user.Name,
		OwnerId:    user.ID,
		BingoCards: bingoCards,
		Users:      []entity.User{*user},
	}

	bingoBoard, err = u.bingoBoardRepo.CreatePersonalBingoBoard(bingoBoard)
	if err != nil {
		return fmt.Errorf("BingoBoardUseCase - CreatePersonalBingoBoard - u.bingoBoardRepo.CreatePersonalBingoBoard: %w", err)
	}

	return nil

}

func (u *BingoBoardUseCase) DeletePersonalBingoBoard(user *entity.User) error {
	err := u.bingoBoardRepo.DeletePersonalBingoBoard(user)
	if err != nil {
		return fmt.Errorf("BingoBoardUseCase - DeletePersonalBingoBoard - u.bingoBoardRepo.DeletePersonalBingoBoard: %w", err)
	}
	return nil
}

func (u *BingoBoardUseCase) JoinBingoBoard(user *entity.User, shortUuid string) error {
	err := u.bingoBoardRepo.AddUserToBingoBoard(user, shortUuid)
	if err != nil {
		return fmt.Errorf("BingoBoardUseCase - JoinBingoBoard - u.bingoBoardRepo.AddUserToBingoBoard: %w", err)
	}
	return nil
}

func (u *BingoBoardUseCase) LeaveBingoBoard(user *entity.User, shortUuid string) error {
	bingoBoard, err := u.bingoBoardRepo.GetBingoBoard(shortUuid)
	if err != nil {
		return fmt.Errorf("BingoBoardUseCase - LeaveBingoBoard - u.bingoBoardRepo.GetBingoBoard: %w", err)
	}

	if bingoBoard.OwnerId == user.ID {
		return fmt.Errorf("BingoBoardUseCase - LeaveBingoBoard - user is owner of bingo board")
	}

	err = u.bingoBoardRepo.RemoveUserFromBingoBoard(user, shortUuid)
	if err != nil {
		return fmt.Errorf("BingoBoardUseCase - LeaveBingoBoard - u.bingoBoardRepo.RemoveUserFromBingoBoard: %w", err)
	}
	return nil
}
