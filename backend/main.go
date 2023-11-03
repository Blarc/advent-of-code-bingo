package main

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	models.ConnectDatabase()

	app := App{
		router: gin.Default(),
	}
	app.start()
}
