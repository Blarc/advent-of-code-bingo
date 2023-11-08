package main

import (
	"github.com/Blarc/advent-of-code-bingo/auth"
	"github.com/Blarc/advent-of-code-bingo/controllers"
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

	githubOAuth := auth.OAuth{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("GITHUB_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("GITHUB_CLIENT_SECRET"),
			Endpoint:     github.Endpoint,
			RedirectURL:  utils.GetEnvVariable("GITHUB_REDIRECT_URI"),
			Scopes:       nil,
		},
	}

	googleOAuth := auth.OAuth{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("GOOGLE_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("GOOGLE_CLIENT_SECRET"),
			Endpoint:     google.Endpoint,
			RedirectURL:  utils.GetEnvVariable("GOOGLE_REDIRECT_URI"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile"},
		},
	}

	redditOAuth := auth.OAuth{
		Config: &oauth2.Config{
			ClientID:     utils.GetEnvVariable("REDDIT_CLIENT_ID"),
			ClientSecret: utils.GetEnvVariable("REDDIT_CLIENT_SECRET"),
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://www.reddit.com/api/v1/authorize",
				TokenURL:  "https://www.reddit.com/api/v1/access_token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
			RedirectURL: utils.GetEnvVariable("REDDIT_REDIRECT_URI"),
			Scopes:      []string{"identity"},
		},
		UserAgent: utils.GetEnvVariable("REDDIT_USER_AGENT"),
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
	apiPublic.GET("/auth/github/callback", func(context *gin.Context) {
		auth.GithubCallbackHandler(context, &githubOAuth)
	})
	apiPublic.GET("/auth/google", googleOAuth.LoginRedirectHandler)
	apiPublic.GET("/auth/google/callback", func(context *gin.Context) {
		auth.GoogleCallbackHandler(context, &googleOAuth)
	})
	apiPublic.GET("/auth/reddit", redditOAuth.LoginRedirectHandler)
	apiPublic.GET("/auth/reddit/callback", func(context *gin.Context) {
		auth.RedditCallbackHandler(context, &redditOAuth)
	})

	// Protected API
	protected := app.router.Group("/api/v1")
	protected.Use(auth.Verifier())
	protected.GET("/me", controllers.FindMe)
	protected.GET("/bingoCards", controllers.FindBingoCards)

	log.Fatal(app.router.Run())
}

func (app *App) Health(c *gin.Context) {
	log.Printf("%v\n", c.Request)
	log.Println(c.Request.Cookies())
	c.String(http.StatusOK, "OK")
}
