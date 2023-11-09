package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BingoCardDto struct {
	ID          uint   `json:"id"`
	Description string `json:"description"`
	UserCount   uint   `json:"user_count"`
}

func FindBingoCards(c *gin.Context) {
	var bingoCards []BingoCardDto

	models.DB.Table("bingo_cards").
		Select("bingo_cards.id, bingo_cards.description, count(user_bingo_cards.bingo_card_id) as user_count").
		Joins("left join user_bingo_cards on user_bingo_cards.bingo_card_id = bingo_cards.id").
		Group("bingo_cards.id").
		Scan(&bingoCards)

	c.JSON(http.StatusOK, bingoCards)
}
