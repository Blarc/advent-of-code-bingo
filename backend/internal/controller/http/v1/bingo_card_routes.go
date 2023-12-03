package v1

import (
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type bingoCardRoutes struct {
	b usecase.BingoCard
	l logger.Interface
}

func newBingoCardRoutesPublic(handler *gin.RouterGroup, b usecase.BingoCard, l logger.Interface) {
	r := &bingoCardRoutes{b, l}

	h := handler.Group("/bingoCard")
	{
		h.GET("/", r.getBingoCards)
	}
}

func newBingoCardRoutesProtected(handler *gin.RouterGroup, b usecase.BingoCard, l logger.Interface) {
	r := &bingoCardRoutes{b, l}

	h := handler.Group("/bingoCard")
	{
		h.POST("/:id/click", r.clickBingoCard)

	}
}

// getBingoCards godoc
// @Summary Get all bingo cards.
// @Description Get all bingo cards.
// @Tags Bingo Card
// @Accept json
// @Produce json
// @Success 200 {array} entity.BingoCardDto
// @Router /bingoCard [get]
func (r *bingoCardRoutes) getBingoCards(c *gin.Context) {
	var bingoCards []entity.BingoCardDto

	bingoCards, err := r.b.GetBingoCards()
	if err != nil {
		r.l.Error(err, "http - v1 - getBingoCards")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, bingoCards)
}

type clickBingoCardRequest struct {
	ID string `uri:"id" binding:"required"`
}

// clickBingoCard godoc
// @Summary Click bingo card.
// @Description Add or remove bingo card from user.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserDto
// @Router /bingoCard/{id}/click [post]
func (r *bingoCardRoutes) clickBingoCard(c *gin.Context) {
	var request clickBingoCardRequest
	if err := c.ShouldBindUri(&request); err != nil {
		r.l.Error(err, "http - v1 - clickBingoCard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	cardUuid, err := uuid.Parse(request.ID)
	if err != nil {
		r.l.Error(err, "http - v1 - clickBingoCard")
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	var user = c.MustGet("user").(*entity.User)
	err = r.b.ClickBingoCard(user, cardUuid)
	if err != nil {
		r.l.Error(err, "http - v1 - clickBingoCard")
		errorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, user.MapToDto())
}
