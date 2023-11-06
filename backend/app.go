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
	"strings"
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
			RedirectURL:  utils.GetEnvVariable("GITHUB_REDIRECT_URI"),
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
			RedirectURL: utils.GetEnvVariable("REDDIT_REDIRECT_URI"),
			Scopes:      []string{"identity"},
		},
	}

	googleOAuth := gin_oauth2.GinOAuth2{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("GOOGLE_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("GOOGLE_CLIENT_SECRET"),
			Endpoint:     google.Endpoint,
			RedirectURL:  utils.GetEnvVariable("GOOGLE_REDIRECT_URI"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		},
	}

	oAuthVerifier := gin_oauth2.GinOAuth2Verifier{
		GithubConfig: githubOAuth.Config,
		GoogleConfig: googleOAuth.Config,
		RedditConfig: redditOAuth.Config,
	}

	app.router.Use(sessions.Sessions("session", sessions.NewCookieStore([]byte("secret"))))

	// Frontend
	app.router.Use(static.Serve("/", static.LocalFile("../frontend/dist/frontend", true)))
	app.router.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.File("../frontend/dist/frontend/index.html")
		}
		// default 404 page not found
	})

	// Public API
	apiPublic := app.router.Group("/api/v1")
	apiPublic.GET("/health", app.Health)
	apiPublic.GET("/auth/github", githubOAuth.LoginRedirectHandler)
	apiPublic.GET("/auth/reddit", redditOAuth.LoginRedirectHandler)
	apiPublic.GET("/auth/google", googleOAuth.LoginRedirectHandler)
	apiPublic.GET("/auth/github/callback", githubOAuth.CallbackHandler)
	apiPublic.GET("/auth/reddit/callback", redditOAuth.CallbackHandler)
	apiPublic.GET("/auth/google/callback", googleOAuth.CallbackHandler)

	// Protected API
	auth := app.router.Group("/api/v1")
	auth.Use(oAuthVerifier.AuthVerifier())
	auth.GET("/me", controllers.FindMe)
	auth.GET("/bingoCards", controllers.FindBingoCards)

	log.Fatal(app.router.Run())
}

func (app *App) Health(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
