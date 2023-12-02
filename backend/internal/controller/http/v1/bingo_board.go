package v1

import (
	models "github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type bingoBoardRoutes struct {
	b usecase.BingoBoard
	l logger.Interface
}

func newBingoBoardRoutes(handler *gin.RouterGroup, b usecase.BingoBoard, l logger.Interface) {
	r := &bingoBoardRoutes{b, l}

	h := handler.Group("/bingoBoard")
	{
		h.GET("/:shortUuid", r.getBingoBoard)
	}
}

type getBingoBoardRequest struct {
	ShortUUID string `uri:"shortUuid" binding:"required"`
}

// GetBingoBoard
// @Summary 	Get bingo board.
// @Description Get a bingo board.
// @Tags 		Bingo Board
// @Accept 		json
// @Produce 	json
// @Success 	200 {object} models.BingoBoardDto
// @Failure 	400 {object} response
// @Router 		/bingoBoard/{id} [get]
// @Param 		id path string true "Bingo Board ID"
// @Security 	Token
func (r *bingoBoardRoutes) getBingoBoard(c *gin.Context) {
	var request getBingoBoardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - getBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	// Find the bingo board by the first 16 characters of the ID
	bingoBoardDto, err := r.b.GetBingoBoard(request.ShortUUID)
	if err != nil {
		r.l.Error(err, "http - v1 - getBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, bingoBoardDto)
}

// CreateBingoBoard godoc
// @Summary Create bingo board.
// @Description Create a new bingo board with random bingo cards.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} models.UserDto
// @Router /bingoBoard [post]
// @Security Token
func (r *bingoBoardRoutes) createBingoBoard(c *gin.Context) {
	var user = c.MustGet("user").(*entity.User)

	board, err := r.b.CreatePersonalBingoBoard(user)
	if err != nil {
		r.l.Error(err, "http - v1 - createBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}
