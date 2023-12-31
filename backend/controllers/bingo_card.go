package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FindBingoCards godoc
// @Summary Get all bingo cards.
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
		Where("bingo_cards.public = true").
		Scan(&bingoCards)

	c.JSON(http.StatusOK, bingoCards)
}

type CreateBingoCardDto struct {
	Description string `json:"description" binding:"required"`
}

// CreateBingoCard godoc
// @Summary Create a bingo card.
// @Description Create a bingo card.
// @Tags Bingo Card
// @Accept json
// @Produce json
// @Param bingoCard body CreateBingoCardDto true "Bingo Card"
// @Success 200 {object} models.BingoCardDto
// @Router /bingoCards [post]
func CreateBingoCard(c *gin.Context) {
	var json CreateBingoCardDto
	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	bingoCard := models.BingoCard{
		Description: json.Description,
		Public:      false,
	}

	models.DB.Create(&bingoCard)
	c.JSON(http.StatusOK, bingoCard.MapToDto())
}
