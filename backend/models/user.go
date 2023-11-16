package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID         uuid.UUID   `gorm:"type:uuid;default:gen_random_uuid()"`
	GithubID   string      `gorm:"index:auth_id,unique" json:"github_id"`
	GoogleID   string      `gorm:"index:auth_id,unique" json:"google_id"`
	RedditID   string      `gorm:"index:auth_id,unique" json:"reddit_id"`
	AvatarURL  string      `gorm:"size:255" json:"avatar_url"`
	Name       string      `gorm:"size:255" json:"name"`
	GithubURL  string      `gorm:"size:255" json:"github_url"`
	RedditURL  string      `gorm:"size:255" json:"reddit_url"`
	BingoCards []BingoCard `gorm:"many2many:user_bingo_card;" json:"bingo_cards"`
}

type UserDto struct {
	Name         string         `json:"name"`
	AvatarURL    string         `json:"avatar_url"`
	GithubURL    string         `json:"github_url"`
	RedditURL    string         `json:"reddit_url"`
	BingoCardDto []BingoCardDto `json:"bingo_cards"`
}

func (u *User) MapToDto() UserDto {
	var bingoCardDtos []BingoCardDto
	for _, bingoCard := range u.BingoCards {
		bingoCardDtos = append(bingoCardDtos, bingoCard.MapToDto())
	}

	return UserDto{
		Name:         u.Name,
		AvatarURL:    u.AvatarURL,
		GithubURL:    u.GithubURL,
		RedditURL:    u.RedditURL,
		BingoCardDto: bingoCardDtos,
	}
}
