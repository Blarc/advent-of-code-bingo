package v1

import (
	"github.com/Blarc/advent-of-code-bingo/internal/entity"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"strings"

	// Swagger docs.
	_ "github.com/Blarc/advent-of-code-bingo/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// FindMe godoc
// @Summary Get user information.
// @Description Get information about the user that is currently logged in.
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} entity.UserDto
// @Failure 401 {object} response
// @Router /me [get]
// @Security Token
func FindMe(c *gin.Context) {
	userDto := c.MustGet("user").(*entity.UserDto)
	c.JSON(http.StatusOK, userDto)
}

// NewRouter -.
// Swagger spec:
// @title       Advent of Code Bingo API
// @version     1.0
// @description Advent of Code Bingo API in Go using Gin framework.
// @host        localhost:8080
// @BasePath    /api/v1
// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func NewRouter(handler *gin.Engine, l logger.Interface, u usecase.User, b usecase.BingoBoard, c usecase.BingoCard, a *AuthMiddleware) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Frontend
	handler.Use(static.Serve("/", static.LocalFile("../frontend/dist/frontend", true)))
	handler.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.File("../frontend/dist/frontend/index.html")
		}
		// default 404 page not found
	})

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/api/v1/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	public := handler.Group("/api/v1")
	{
		newAuthRoutes(public, u, l, a)
		newBingoCardRoutesPublic(public, c, l)
	}

	protected := handler.Group("/api/v1")
	{
		protected.Use(a.Verifier())
		protected.GET("/me", FindMe)
		newBingoBoardRoutes(protected, b, l, a)
		newBingoCardRoutesProtected(protected, c, l)
	}
}
