package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	OAuthID    uint        `gorm:"unique" json:"oauth_id"`
	AvatarURL  string      `gorm:"size:255" json:"avatar_url"`
	Name       string      `gorm:"size:255" json:"name"`
	GitHubURL  string      `gorm:"size:255" json:"github_url"`
	BingoCards []BingoCard `gorm:"many2many:user_bingo_cards;" json:"bingo_cards"`
}
