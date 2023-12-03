package v1

import (
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type bingoBoardRoutes struct {
	b usecase.BingoBoard
	l logger.Interface
}

func newBingoBoardRoutes(handler *gin.RouterGroup, b usecase.BingoBoard, l logger.Interface, a *AuthMiddleware) {
	r := &bingoBoardRoutes{b, l}

	h := handler.Group("/bingoBoard").Use(a.Verifier())
	{
		h.GET("/:id", r.getBingoBoard)
		h.POST("/", r.createBingoBoard)
		h.DELETE("/", r.deleteBingoBoard)
		h.POST("/:id/join", r.JoinBingoBoard)
		h.DELETE("/:id/leave", r.LeaveBingoBoard)
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
// @Success 	200 {object} entity.BingoBoardDto
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
// @Success 200 {object} entity.UserDto
// @Router /bingoBoard [post]
// @Security Token
func (r *bingoBoardRoutes) createBingoBoard(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	err := r.b.CreatePersonalBingoBoard(user)
	if err != nil {
		r.l.Error(err, "http - v1 - createBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

// DeleteBingoBoard godoc
// @Summary Delete bingo board.
// @Description Irrevocably delete a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserDto
// @Router /bingoBoard/{id} [delete]
// @Security Token
func (r *bingoBoardRoutes) deleteBingoBoard(c *gin.Context) {
	user := c.MustGet("user").(*entity.User)

	err := r.b.DeletePersonalBingoBoard(user)
	if err != nil {
		r.l.Error(err, "http - v1 - deleteBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

type JoinBingoBoardRequest struct {
	ShortUUID string `uri:"shortUuid" binding:"required"`
}

// JoinBingoBoard godoc
// @Summary Join bingo board.
// @Description Join a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserDto
// @Router /bingoBoard/{id}/join [post]
// @Param id path string true "Bingo Board ID"
// @Security Token
func (r *bingoBoardRoutes) JoinBingoBoard(c *gin.Context) {
	var request JoinBingoBoardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - JoinBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(*entity.User)
	err := r.b.JoinBingoBoard(user, request.ShortUUID)
	if err != nil {
		r.l.Error(err, "http - v1 - JoinBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}

type LeaveBingoBoardRequest struct {
	ShortUUID string `uri:"shortUuid" binding:"required"`
}

// LeaveBingoBoard godoc
// @Summary Leave bingo board.
// @Description Leave a bingo board.
// @Tags Bingo Board
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserDto
// @Router /bingoBoard/{id}/leave [delete]
// @Param id path string true "Bingo Board ID"
// @Security Token
func (r *bingoBoardRoutes) LeaveBingoBoard(c *gin.Context) {
	var request LeaveBingoBoardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - LeaveBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	user := c.MustGet("user").(*entity.User)
	err := r.b.LeaveBingoBoard(user, request.ShortUUID)
	if err != nil {
		r.l.Error(err, "http - v1 - LeaveBingoBoard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}
