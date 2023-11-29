package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"log"
	"math/rand"
	"net/http"
)

// CreateBingoBoard godoc
// @Summary Create bingo board.
// @Schemes http
// @Description Create a new bingo board with random bingo cards.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /bingoBoard [post]
// @Security Token
func CreateBingoBoard(c *gin.Context) {
	var user = c.MustGet("user").(models.User)
	if user.PersonalBingoBoard != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You already have a personal bingo board"})
		return
	}

	// Find all bingo cards
	var bingoCards []models.BingoCard
	models.DB.Model(&models.BingoCard{}).Find(&bingoCards)
	// Shuffle them
	rand.Shuffle(len(bingoCards), func(i, j int) {
		bingoCards[i], bingoCards[j] = bingoCards[j], bingoCards[i]
	})

	// Create the new bingo board
	newBoard := models.BingoBoard{
		Name:    user.Name,
		OwnerId: user.ID,
		// Take the first 16 cards from the shuffled list
		BingoCards: bingoCards[:16],
	}
	models.DB.Create(&newBoard)

	// Add the new bingo board to the user
	err := models.DB.Model(&user).Association("PersonalBingoBoard").Append(&newBoard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

// FindBingoBoard godoc
// @Summary Get bingo board.
// @Schemes http
// @Description Get a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /bingoBoard/{id} [get]
// @Param id path string true "Bingo Board ID"
// @Security Token
func FindBingoBoard(c *gin.Context) {
	var bingoBoardId models.BingoBoardId
	if err := c.ShouldBindUri(&bingoBoardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the bingo board by the first 16 characters of the ID
	var bingoBoard models.BingoBoard
	result := models.DB.
		First(&bingoBoard, "substring(id::text, 1, 16) = ?", bingoBoardId.ID)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	// Find all bingo cards for the board and count how many users have them
	var bingoCards []models.BingoCardDto
	result = models.DB.Table("bingo_cards").
		Select("bingo_cards.id, bingo_cards.description, count(user_bingo_card.user_id) as user_count").
		Joins("left join bingo_board_bingo_card ON bingo_board_bingo_card.bingo_card_id = bingo_cards.id").
		Joins("left join user_bingo_card ON user_bingo_card.bingo_card_id = bingo_cards.id").
		Group("bingo_board_bingo_card.bingo_board_id, bingo_cards.id").
		Find(&bingoCards, "substring(bingo_board_bingo_card.bingo_board_id::text, 1, 16) = ?", bingoBoardId.ID)

	log.Printf("%v\n", bingoCards)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	// Create the DTO
	bingoBoardDto := bingoBoard.MapToDto()
	bingoBoardDto.BingoCards = bingoCards
	c.JSON(http.StatusOK, bingoBoardDto)
}

// DeleteBingoBoard godoc
// @Summary Delete bingo board.
// @Schemes http
// @Description Irrevocably delete a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /bingoBoard/{id} [delete]
// @Security Token
func DeleteBingoBoard(c *gin.Context) {

	var user = c.MustGet("user").(models.User)
	if user.PersonalBingoBoard == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You don't have a personal bingo board"})
		return
	}

	err := models.DB.Unscoped().Model(&user.PersonalBingoBoard).Association("BingoCards").Clear()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = models.DB.Model(&user.PersonalBingoBoard).Association("Users").Clear()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Unscoped().Delete(&user.PersonalBingoBoard)
	user.PersonalBingoBoard = nil
	c.JSON(http.StatusOK, user.MapToDto())
}

// JoinBingoBoard godoc
// @Summary Join bingo board.
// @Schemes http
// @Description Join a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /bingoBoard/{id}/join [post]
// @Param id path string true "Bingo Board ID"
// @Security Token
func JoinBingoBoard(c *gin.Context) {
	var bingoBoardId models.BingoBoardId
	if err := c.ShouldBindUri(&bingoBoardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the bingo board by the first 16 characters of the ID
	var bingoBoard models.BingoBoard
	if err := models.DB.First(&bingoBoard, "substring(id::text, 1, 16) = ?", bingoBoardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user = c.MustGet("user").(models.User)
	err := models.DB.Model(&user).Association("BingoBoards").Append(&bingoBoard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

// LeaveBingoBoard godoc
// @Summary Leave bingo board.
// @Schemes http
// @Description Leave a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /bingoBoard/{id}/leave [delete]
// @Param id path string true "Bingo Board ID"
// @Security Token
func LeaveBingoBoard(c *gin.Context) {
	var bingoBoardId models.BingoBoardId
	if err := c.ShouldBindUri(&bingoBoardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the bingo board by the first 16 characters of the ID
	var bingoBoard models.BingoBoard
	if err := models.DB.First(&bingoBoard, "substring(id::text, 1, 16) = ?", bingoBoardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user = c.MustGet("user").(models.User)
	if user.ID == bingoBoard.OwnerId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot leave your own board"})
		return
	}

	err := models.DB.Model(&user).Association("BingoBoards").Delete(&bingoBoard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

// AddBingoCard godoc
// @Summary Add bingo card.
// @Schemes http
// @Description Add a bingo card to a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /bingoBoard/{id}/addBingoCard [put]
// @Param id path string true "Bingo Board ID"
// @Param data body models.BingoCardId true "Bingo Card UUID"
// @Security Token
func AddBingoCard(c *gin.Context) {
	var bingoBoardId models.BingoBoardId
	if err := c.ShouldBindUri(&bingoBoardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the bingo board by the first 16 characters of the ID
	var bingoBoard models.BingoBoard
	if err := models.DB.First(&bingoBoard, "substring(id::text, 1, 16) = ?", bingoBoardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bingoCardId models.BingoCardId
	if err := c.ShouldBindJSON(&bingoCardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bingoCard models.BingoCard
	if err := models.DB.First(&bingoCard, bingoCardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.DB.Model(&bingoBoard).Association("BingoCards").Append(&bingoCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bingoBoard.MapToDto())
}

// RemoveBingoCard godoc
// @Summary Remove bingo card.
// @Schemes http
// @Description Remove a bingo card from a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} string
// @Router /bingoBoard/{id}/removeBingoCard [put]
// @Param id path string true "Bingo Board ID"
// @Param data body models.BingoCardId true "Bingo Card UUID"
// @Security Token
func RemoveBingoCard(c *gin.Context) {
	var bingoBoardId models.BingoBoardId
	if err := c.ShouldBindUri(&bingoBoardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find the bingo board by the first 16 characters of the ID
	var bingoBoard models.BingoBoard
	if err := models.DB.First(&bingoBoard, "substring(id::text, 1, 16) = ?", bingoBoardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bingoCardId models.BingoCardId
	if err := c.ShouldBindJSON(&bingoCardId); err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var bingoCard models.BingoCard
	if err := models.DB.First(&bingoCard, bingoCardId.ID).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := models.DB.Model(&bingoBoard).Association("BingoCards").Delete(&bingoCard)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, bingoBoard.MapToDto())
}
