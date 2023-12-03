package app

import (
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/config"
	v1 "github.com/Blarc/advent-of-code-bingo/internal/controller/http/v1"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase/repo"
	"github.com/Blarc/advent-of-code-bingo/pkg/httpserver"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/Blarc/advent-of-code-bingo/pkg/postgres"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
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

	bingoBoardUseCase := usecase.NewBingoBoardUseCase(
		repo.NewBingoBoardRepo(pg),
		repo.NewBingoCardRepo(pg),
	)

	bingoCardUseCase := usecase.NewBingoCardUseCase(
		repo.NewBingoCardRepo(pg),
	)

	userUseCase := usecase.NewUserUseCase(
		repo.NewUserRepo(pg),
	)

	githubOAuth := &oauth2.Config{
		ClientID:     cfg.OAuth.GithubClientId,
		ClientSecret: cfg.OAuth.GithubClientSecret,
		Endpoint:     github.Endpoint,
		RedirectURL:  cfg.OAuth.GithubRedirectUri,
		Scopes:       nil,
	}

	googleOAuth := &oauth2.Config{
		ClientID:     cfg.OAuth.GithubClientId,
		ClientSecret: cfg.OAuth.GithubClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  cfg.OAuth.GithubRedirectUri,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
	}

	redditOAuth := &oauth2.Config{
		ClientID:     cfg.OAuth.RedditClientId,
		ClientSecret: cfg.OAuth.RedditClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   "https://www.reddit.com/api/v1/authorize",
			TokenURL:  "https://www.reddit.com/api/v1/access_token",
			AuthStyle: oauth2.AuthStyleInHeader,
		},
		RedirectURL: cfg.OAuth.RedditRedirectUri,
		Scopes:      []string{"identity"},
	}

	handler := gin.New()

	authMiddleware := v1.NewAuthMiddleware(userUseCase, l, cfg.OAuth.TokenEncryptSecret, githubOAuth, googleOAuth, redditOAuth)
	v1.NewRouter(handler, l, userUseCase, bingoBoardUseCase, bingoCardUseCase, authMiddleware)
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
