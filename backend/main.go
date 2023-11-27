package main

import (
	"github.com/Blarc/advent-of-code-bingo/docs"
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/Blarc/advent-of-code-bingo/utils"
	"github.com/joho/godotenv"
	"log"
	"os"
)

// @title           Advent of Code Bingo API
// @version         1.0
// @description     Advent of Code Bingo API in Go using Gin framework.

// @contact.name Jakob Maležič
// @contact.url https://github.com/Blarc

// @license.name  GNU General Public License v3.0
// @license.url   https://www.gnu.org/licenses/gpl-3.0.html

// @BasePath  /api/v1
// @securityDefinitions.apikey Token
// @in header
// @name Authorization
func main() {
	log.SetOutput(os.Stdout)

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found.")
	}

	docs.SwaggerInfo.Host = utils.GetEnvVariable("HOST")
	docs.SwaggerInfo.Schemes = []string{utils.GetEnvVariable("SCHEME")}
	models.ConnectDatabase()

	newApp().start()
}
