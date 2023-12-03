package usecase

import (
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/google/go-github/v56/github"
	"github.com/google/uuid"
	"google.golang.org/api/oauth2/v2"
)

//go:generate mockgen -source=interfaces.go -destination=mocks/mocks_test.go -package=usecase_test

type (
	BingoBoard interface {
		GetBingoBoard(shortUuid string) (*entity.BingoBoardDto, error)
		CreatePersonalBingoBoard(user *entity.User) error
		DeletePersonalBingoBoard(user *entity.User) error
		JoinBingoBoard(user *entity.User, shortUuid string) error
		LeaveBingoBoard(user *entity.User, shortUuid string) error
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
		ClickBingoCard(user *entity.User, id uuid.UUID) error
	}

	BingoCardRepo interface {
		GetBingoCard(id uuid.UUID) (entity.BingoCard, error)
		GetBingoCards() ([]entity.BingoCard, error)
		GetBingoCardsWithUserCount() ([]entity.BingoCardDto, error)
		AddBingoCardToUser(user *entity.User, bingoCard *entity.BingoCard) error
		RemoveBingoCardFromUser(user *entity.User, bingoCard *entity.BingoCard) error
	}

	User interface {
		GetUser(id uuid.UUID) (*entity.UserDto, error)
		CreateGithubUser(githubUser *github.User) (*uuid.UUID, error)
		CreateGoogleUser(googleUser *oauth2.Userinfo) (*uuid.UUID, error)
		CreateRedditUser(redditUser *RedditUserData) (*uuid.UUID, error)
	}

	UserRepo interface {
		GetUser(id uuid.UUID) (*entity.User, error)
		CreateGithubUser(githubUser *github.User) (*entity.User, error)
		CreateGoogleUser(googleUser *oauth2.Userinfo) (*entity.User, error)
		CreateRedditUser(redditUser *RedditUserData) (*entity.User, error)
	}

	// Could add additional interfaces here for other use cases
	// For example for calling other microservices (APIs)
)
