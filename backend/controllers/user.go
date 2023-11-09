package controllers

import (
	"encoding/gob"
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func init() {
	gob.Register(models.User{})
}

func FindMe(c *gin.Context) {
	c.JSON(http.StatusOK, c.MustGet("user"))
}

type BingoCardId struct {
	ID uint `uri:"id" binding:"required"`
}

func ClickBingoCard(c *gin.Context) {

	var bingoCardId BingoCardId
	if err := c.ShouldBindUri(&bingoCardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bingoCard models.BingoCard
	if err := models.DB.First(&bingoCard, bingoCardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user = c.MustGet("user").(models.User)
	for _, bingoCard := range user.BingoCards {
		if bingoCard.ID == bingoCardId.ID {
			log.Printf("bingo card %d already clicked", bingoCardId.ID)
			err := models.DB.Model(&user).Association("BingoCards").Delete(&bingoCard)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			}
			c.JSON(http.StatusOK, user)
			return
		}
	}

	err := models.DB.Model(&user).Association("BingoCards").Append(&bingoCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
