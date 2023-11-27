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

// FindMe godoc
// @Summary Get user information.
// @Schemes http
// @Description Get information about the user that is currently logged in.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /me [get]
// @Security Token
func FindMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	log.Printf("user: %+v", user)
	c.JSON(http.StatusOK, user.MapToDto())
}

// ClickBingoCard godoc
// @Summary Click bingo card.
// @Schemes http
// @Description Add or remove bingo card from user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /me/bingoCard/{id} [post]
// @Param id path int true "Bingo Card ID"
// @Security Token
func ClickBingoCard(c *gin.Context) {

	var bingoCardId models.BingoCardId
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

	c.JSON(http.StatusOK, user.MapToDto())
}
