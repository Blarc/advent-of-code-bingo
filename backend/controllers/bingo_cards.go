package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/auth"
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type BingoCardDto struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	UserCount   uint   `json:"user_count"`
	Selected    bool   `json:"selected"`
}

func FindBingoCards(c *gin.Context) {
	var bingoCards []BingoCardDto

	models.DB.Table("bingo_cards").
		Select("bingo_cards.id, bingo_cards.description, count(user_bingo_cards.bingo_card_id) as user_count").
		Joins("left join user_bingo_cards on user_bingo_cards.bingo_card_id = bingo_cards.id").
		Group("bingo_cards.id").
		Scan(&bingoCards)

	userUuid := auth.GetUserUuidFromHeader(c)
	log.Println(userUuid)
	if userUuid != nil {

		var user models.User
		result := models.DB.Preload("BingoCards").First(&user, "id = ?", userUuid)
		if result.Error != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// TODO: Optimize this
		for i, bingoCard := range bingoCards {
			for _, selectedBingoCard := range user.BingoCards {
				if bingoCard.ID == selectedBingoCard.ID {
					bingoCards[i].Selected = true
				}
			}
		}
	}

	c.JSON(http.StatusOK, bingoCards)
}
