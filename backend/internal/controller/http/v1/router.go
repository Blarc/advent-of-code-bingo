package v1

import (
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
)

// NewRouter -.
// Swagger spec:
// @title       Advent of Code Bingo API
// @version     1.0
// @description Advent of Code Bingo API in Go using Gin framework.
// @host        localhost:8080
// @BasePath    /v1
func NewRouter(handler *gin.Engine, l logger.Interface, t usecase.BingoBoard) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// K8s probe
	handler.GET("/healthz", func(c *gin.Context) { c.Status(http.StatusOK) })

	// Routers
	h := handler.Group("/v1")
	{
		newBingoBoardRoutes(h, t, l)
	}
}
