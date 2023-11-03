package main

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type App struct {
	router *gin.Engine
}

func (app *App) start() {
	// Frontend
	app.router.Use(static.Serve("/", static.LocalFile("../frontend/dist/frontend", false)))
	// Backend
	app.router.GET("/api/v1/health", app.health)
	log.Fatal(app.router.Run())
}

func (app *App) health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
