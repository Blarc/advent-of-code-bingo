package controllers

import (
	"encoding/gob"
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"log"
	"net/http"
)

func init() {
	gob.Register(models.User{})
}

func LogInUserGitHub(c *gin.Context, conf *oauth2.Config) {

	accessToken := c.MustGet("token").(oauth2.Token)
	githubClient := github.NewClient(conf.Client(c, &accessToken))
	githubUserData, _, err := githubClient.Users.Get(c, "")
	if err != nil {
		return
	}

	var user models.User
	models.DB.FirstOrCreate(&user, models.User{
		OAuthID:   *githubUserData.ID,
		Name:      *githubUserData.Name,
		AvatarURL: *githubUserData.AvatarURL,
		GitHubURL: *githubUserData.HTMLURL,
	})

	log.Printf("Saving user to session: %v", user)
	session := sessions.Default(c)
	session.Set("user", user)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func LogInUserReddit(c *gin.Context) {
	//accessToken := c.MustGet("token").(oauth2.Token)
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func LogInUserGoogle(c *gin.Context, conf *oauth2.Config) {
	accessToken := c.MustGet("token").(oauth2.Token)

	googleService, err := goauth.NewService(c, option.WithTokenSource(conf.TokenSource(c, &accessToken)))
	if err != nil {
		log.Printf("Failed to create oauth service: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	userInfo, err := googleService.Userinfo.Get().Do()
	if err != nil {
		log.Printf("Failed to get userinfo for user: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var user models.User
	models.DB.FirstOrCreate(&user, models.User{
		Name:      userInfo.Name,
		AvatarURL: userInfo.Picture,
	})

	log.Printf("Saving user to session: %v", user)
	session := sessions.Default(c)
	session.Set("user", user)
	if err := session.Save(); err != nil {
		log.Printf("Failed to save session: %v", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusMovedPermanently, "/login")
}

func FindMe(c *gin.Context) {
	c.JSON(http.StatusOK, sessions.Default(c).Get("user"))
}
