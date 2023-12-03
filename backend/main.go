package main

import (
	"github.com/Blarc/advent-of-code-bingo/config"
	"github.com/Blarc/advent-of-code-bingo/internal/app"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	app.Run(cfg)

}
