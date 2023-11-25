package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/auth"
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// FindBingoCards godoc
// @Summary Get all bingo cards.
// @Schemes http
// @Description Get all bingo cards.
// @Tags Bingo Card
// @Accept json
// @Produce json
// @Success 200 {array} models.BingoCardDto
// @Router /bingoCards [get]
func FindBingoCards(c *gin.Context) {
	var bingoCards []models.BingoCardDto

	models.DB.Table("bingo_cards").
		Select("bingo_cards.id, bingo_cards.description, count(user_bingo_card.bingo_card_id) as user_count").
		Joins("left join user_bingo_card on user_bingo_card.bingo_card_id = bingo_cards.id").
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
