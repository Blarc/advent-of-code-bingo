package models

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	dbHost := utils.GetEnvVariable("DB_HOST")
	dbUser := utils.GetEnvVariable("DB_USER")
	dbPass := utils.GetEnvVariable("DB_PASS")
	dbName := utils.GetEnvVariable("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s TimeZone=Europe/Ljubljana",
		dbHost,
		dbUser,
		dbPass,
		dbName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database!")
	}

	err = db.AutoMigrate(&User{}, &BingoCard{})
	if err != nil {
		return
	}

	DB = db
}
