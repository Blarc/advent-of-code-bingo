package usecase

import "github.com/Blarc/advent-of-code-bingo/internal/entity"

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks_test.go -package=usecase_test

type (
	BingoBoard interface {
		GetBingoBoard(shortUuid string) (*entity.BingoBoardDto, error)
		CreatePersonalBingoBoard(user *entity.User) (*entity.BingoBoardDto, error)
		DeletePersonalBingoBoard(user *entity.User) error
		JoinBingoBoard(user *entity.User, shortUuid string) (*entity.User, error)
		LeaveBingoBoard(user *entity.User, shortUuid string) (*entity.User, error)
	}

	BingoBoardRepo interface {
		GetBingoBoard(shortUuid string) (*entity.BingoBoard, error)
		GetBingoCardsWithCount(shortUuid string) ([]entity.BingoCardDto, error)
		CreatePersonalBingoBoard(bingoBoard *entity.BingoBoard) (*entity.BingoBoard, error)
		DeletePersonalBingoBoard(user *entity.User) error
		AddUserToBingoBoard(user *entity.User, shortUuid string) error
		RemoveUserFromBingoBoard(user *entity.User, shortUuid string) error
	}

	BingoCard interface {
		GetBingoCards() ([]entity.BingoCardDto, error)
	}

	BingoCardRepo interface {
		GetBingoCards() ([]entity.BingoCard, error)
		GetBingoCardsWithUserCount() ([]entity.BingoCardDto, error)
	}

	// Could add additional interfaces here for other use cases
	// For example for calling other microservices (APIs)
)
