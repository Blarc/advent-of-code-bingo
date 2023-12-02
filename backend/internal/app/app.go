package app

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/config"
	"github.com/Blarc/advent-of-code-bingo/pkg/httpserver"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/Blarc/advent-of-code-bingo/pkg/postgres"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"syscall"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	pg, err := postgres.New(
		cfg.PG.PgHost,
		cfg.PG.PgPort,
		cfg.PG.PgUser,
		cfg.PG.PgPass,
		cfg.PG.PgName,
	)
	if err != nil {
		l.Fatal(fmt.Errorf("app - Run - postgres.New: %w", err))
	}
	defer func(pg *postgres.Postgres) {
		err := pg.Close()
		if err != nil {
			l.Fatal(fmt.Errorf("app - Run - postgres.Close: %w", err))
		}
	}(pg)

	handler := gin.New()
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: " + s.String())
	case err = <-httpServer.Notify():
		l.Error(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
