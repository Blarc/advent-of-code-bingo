package repo

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/postgres"
	"github.com/google/go-github/v56/github"
	"github.com/google/uuid"
	"google.golang.org/api/oauth2/v2"
	"strconv"
)

type UserRepo struct {
	*postgres.Postgres
}

func NewUserRepo(pg *postgres.Postgres) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) GetUser(userUuid uuid.UUID) (*entity.User, error) {
	var user *entity.User
	result := u.DB.
		Preload("BingoCards").
		Preload("BingoBoards").
		Preload("PersonalBingoBoard").
		First(&user, "id = ?", userUuid)
	if result.Error != nil {
		return nil, fmt.Errorf("UserRepo - GetUser - u.DB: %w", result.Error)
	}
	return user, nil
}

func (u *UserRepo) CreateGithubUser(githubUser *github.User) (*entity.User, error) {
	githubId := strconv.FormatInt(*githubUser.ID, 10)

	var user *entity.User
	result := u.DB.
		Where(entity.User{GithubID: githubId}).
		Assign(entity.User{
			GithubID:  githubId,
			Name:      *githubUser.Login,
			AvatarURL: *githubUser.AvatarURL,
			GithubURL: *githubUser.HTMLURL,
		}).
		FirstOrCreate(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("UserRepo - CreateUser - u.DB: %w", result.Error)
	}
	return user, nil
}

func (u *UserRepo) CreateGoogleUser(googleUser *oauth2.Userinfo) (*entity.User, error) {
	var user *entity.User
	result := u.DB.
		Where(entity.User{GoogleID: googleUser.Id}).
		Assign(entity.User{
			GoogleID:  googleUser.Id,
			Name:      googleUser.Name,
			AvatarURL: googleUser.Picture,
		}).
		FirstOrCreate(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("UserRepo - CreateUser - u.DB: %w", result.Error)
	}
	return user, nil
}

func (u *UserRepo) CreateRedditUser(redditUser *usecase.RedditUserData) (*entity.User, error) {
	var user *entity.User
	result := u.DB.
		Where(entity.User{RedditID: redditUser.ID}).
		Assign(entity.User{
			RedditID:  redditUser.ID,
			Name:      redditUser.Name,
			AvatarURL: redditUser.IconImg,
		}).
		FirstOrCreate(&user)

	if result.Error != nil {
		return nil, fmt.Errorf("UserRepo - CreateUser - u.DB: %w", result.Error)
	}
	return user, nil
}
