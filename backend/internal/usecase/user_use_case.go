package usecase

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/google/go-github/v56/github"
	"github.com/google/uuid"
	"google.golang.org/api/oauth2/v2"
)

type RedditSubredditData struct {
	Url string `json:"url"`
}

type RedditUserData struct {
	ID        string              `json:"id"`
	Name      string              `json:"name"`
	IconImg   string              `json:"icon_img"`
	Subreddit RedditSubredditData `json:"subreddit"`
}

type UserUseCase struct {
	userRepo UserRepo
}

func NewUserUseCase(userRepo UserRepo) *UserUseCase {
	return &UserUseCase{userRepo}
}

func (u *UserUseCase) GetUser(userUuid uuid.UUID) (*entity.UserDto, error) {
	user, err := u.userRepo.GetUser(userUuid)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - GetUser - u.userRepo.GetUser: %w", err)
	}
	userDto := user.MapToDto()
	return &userDto, nil
}

func (u *UserUseCase) CreateGithubUser(githubUser *github.User) (*uuid.UUID, error) {
	user, err := u.userRepo.CreateGithubUser(githubUser)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - CreateGithubUser - u.userRepo.CreateGithubUser: %w", err)
	}
	return &user.ID, nil
}

func (u *UserUseCase) CreateGoogleUser(googleUser *oauth2.Userinfo) (*uuid.UUID, error) {
	user, err := u.userRepo.CreateGoogleUser(googleUser)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - CreateGoogleUser - u.userRepo.CreateGoogleUser: %w", err)
	}
	return &user.ID, nil
}

func (u *UserUseCase) CreateRedditUser(redditUser *RedditUserData) (*uuid.UUID, error) {
	user, err := u.userRepo.CreateRedditUser(redditUser)
	if err != nil {
		return nil, fmt.Errorf("UserUseCase - CreateRedditUser - u.userRepo.CreateRedditUser: %w", err)
	}
	return &user.ID, nil
}
