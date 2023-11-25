package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type CreateBingoBoardDto struct {
	Name string `json:"name"`
}

// CreateBingoBoard godoc
// @Summary Create bingo board.
// @Schemes http
// @Description Create a new bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /me/bingoBoard [post]
// @Param data body CreateBingoBoardDto true "Bingo Board Name"
// @Security Token
func CreateBingoBoard(c *gin.Context) {

	var createBingoBoardDto CreateBingoBoardDto
	if err := c.ShouldBindJSON(&createBingoBoardDto); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user = c.MustGet("user").(models.User)
	err := models.DB.Model(&user).Association("BingoBoards").Append(&models.BingoBoard{
		Name: createBingoBoardDto.Name,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

type BingoBoardId struct {
	ID uint `uri:"id" binding:"required"`
}

// DeleteBingoBoard godoc
// @Summary Delete bingo board.
// @Schemes http
// @Description Delete a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /me/bingoBoard/{id} [delete]
// @Param id path int true "Bingo Board ID"
// @Security Token
func DeleteBingoBoard(c *gin.Context) {

	var bingoBoardId BingoBoardId
	if err := c.ShouldBindUri(&bingoBoardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bingoBoard models.BingoBoard
	if err := models.DB.First(&bingoBoard, bingoBoardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user = c.MustGet("user").(models.User)
	err := models.DB.Model(&user).Association("BingoBoards").Delete(&bingoBoard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Delete(&bingoBoard)

	c.JSON(http.StatusOK, user.MapToDto())
}
