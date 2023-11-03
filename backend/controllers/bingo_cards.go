package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FindBingoCards(c *gin.Context) {
	var bingoCards []models.BingoCard
	models.DB.Find(&bingoCards)

	c.JSON(http.StatusOK, bingoCards)
}
