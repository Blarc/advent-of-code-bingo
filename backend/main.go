package main

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)

	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found.")
	}

	models.ConnectDatabase()

	app := App{
		router: gin.Default(),
	}
	app.start()
}
