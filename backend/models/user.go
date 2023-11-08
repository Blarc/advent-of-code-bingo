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
	BingoCards []BingoCard `gorm:"many2many:user_bingo_cards;" json:"bingo_cards"`
}
