package controllers

import (
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/Blarc/advent-of-code-bingo/utils"
)

func LogInUserGitHub(c *gin.Context) {

	code := c.Query("code")
	accessToken := utils.GetGitHubAccessToken(code)
	gitHubUserData := utils.GetGitHubUserData(accessToken)

	var user models.User
	models.DB.FirstOrCreate(&user, models.User{
		OAuthID:   gitHubUserData.ID,
		Name:      gitHubUserData.Name,
		AvatarURL: gitHubUserData.AvatarURL,
		GitHubURL: gitHubUserData.GitHubURL,
	})

	c.JSON(http.StatusOK, user)
}
