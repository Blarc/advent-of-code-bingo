package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"net/http"

	"github.com/Blarc/advent-of-code-bingo/utils"
)

func LogInUserGitHub(c *gin.Context) {

	accessToken := c.MustGet("token").(oauth2.Token)
	gitHubUserData := utils.GetGitHubUserData(accessToken.AccessToken)

	var user models.User
	models.DB.FirstOrCreate(&user, models.User{
		OAuthID:   gitHubUserData.ID,
		Name:      gitHubUserData.Name,
		AvatarURL: gitHubUserData.AvatarURL,
		GitHubURL: gitHubUserData.GitHubURL,
	})

	c.JSON(http.StatusOK, user)
}

func LogInUserReddit(c *gin.Context) {
	accessToken := c.MustGet("token").(oauth2.Token)
	c.JSON(http.StatusOK, accessToken)
}

func LogInUserGoogle(c *gin.Context) {
	accessToken := c.MustGet("token").(oauth2.Token)
	c.JSON(http.StatusOK, accessToken)
}
