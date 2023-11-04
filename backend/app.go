package main

import (
	"github.com/Blarc/advent-of-code-bingo/controllers"
	"github.com/Blarc/advent-of-code-bingo/gin-oauth2"
	"github.com/Blarc/advent-of-code-bingo/utils"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
)

type App struct {
	router *gin.Engine
}

func (app *App) start() {

	githubOAuth := gin_oauth2.GinOAuth2{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("GITHUB_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("GITHUB_CLIENT_SECRET"),
			Endpoint:     github.Endpoint,
			RedirectURL:  utils.GetEnvVariable("GITHUB_REDIRECT_URL"),
			Scopes:       nil,
		},
	}

	redditOAuth := gin_oauth2.GinOAuth2{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("REDDIT_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("REDDIT_CLIENT_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://www.reddit.com/api/v1/authorize",
				TokenURL: "https://www.reddit.com/api/v1/access_token",
			},
			RedirectURL: utils.GetEnvVariable("REDDIT_REDIRECT_URL"),
			Scopes:      []string{"identity"},
		},
	}

	googleOAuth := gin_oauth2.GinOAuth2{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("GOOGLE_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("GOOGLE_CLIENT_SECRET"),
			Endpoint:     google.Endpoint,
			RedirectURL:  utils.GetEnvVariable("GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
	}

	app.router.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("secret"))))

	// Frontend
	app.router.Use(static.Serve("/", static.LocalFile("../frontend/dist/frontend", false)))

	// Public API
	apiPublic := app.router.Group("/api/v1")
	apiPublic.GET("/health", app.Health)
	apiPublic.GET("/auth/github", func(context *gin.Context) {
		githubOAuth.LoginHandler(context, "GitHub")
	})
	apiPublic.GET("/auth/reddit", func(context *gin.Context) {
		redditOAuth.LoginHandler(context, "Reddit")
	})
	apiPublic.GET("/auth/google", func(context *gin.Context) {
		googleOAuth.LoginHandler(context, "Google")
	})

	// Protected API
	githubAuth := app.router.Group("/api/v1")
	githubAuth.Use(githubOAuth.Auth())
	githubAuth.GET("/auth/github/callback", controllers.LogInUserGitHub)
	// TODO: This goes through, because it only checks if token exists in session cookie
	githubAuth.GET("/bingoCards", controllers.FindBingoCards)

	redditAuth := app.router.Group("/api/v1")
	redditAuth.Use(redditOAuth.Auth())
	redditAuth.GET("/auth/reddit/callback", controllers.LogInUserReddit)

	googleAuth := app.router.Group("/api/v1")
	googleAuth.Use(googleOAuth.Auth())
	googleAuth.GET("/auth/google/callback", controllers.LogInUserGoogle)

	log.Fatal(app.router.Run())
}

func (app *App) Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
