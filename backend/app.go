package main

import (
	"github.com/Blarc/advent-of-code-bingo/controllers"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/zalando/gin-oauth2/github"
	"log"
	"net/http"
)

type App struct {
	router *gin.Engine
}

func (app *App) start() {

	github.Setup(
		"http://localhost:8080/api/v1/login/github/callback",
		"github.json",
		[]string{},
		[]byte("secret"),
	)

	app.router.Use(github.Session("session"))

	// Frontend
	app.router.Use(static.Serve("/", static.LocalFile("../frontend/dist/frontend", false)))

	// Public API
	apiPublic := app.router.Group("/api/v1")
	apiPublic.GET("/health", app.Health)
	apiPublic.GET("/login", github.LoginHandler)

	// Protected API
	apiProtected := app.router.Group("/api/v1")
	apiProtected.Use(github.Auth())
	apiProtected.GET("/login/github/callback", controllers.LogInUserGitHub)
	apiProtected.GET("/bingoCards", controllers.FindBingoCards)

	log.Fatal(app.router.Run())
}

func (app *App) Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func UserInfoHandler(ctx *gin.Context) {
	var (
		res github.AuthUser
		val interface{}
		ok  bool
	)

	log.Println(ctx.Keys)

	val = ctx.MustGet("user")
	if res, ok = val.(github.AuthUser); !ok {
		res = github.AuthUser{
			Name: "no User",
		}
	}

	code := ctx.Query("code")

	ctx.JSON(http.StatusOK, gin.H{"Hello": "from private", "user": res, "code": code})
}
