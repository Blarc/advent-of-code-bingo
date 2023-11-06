package gin_oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"

	goauth "google.golang.org/api/oauth2/v2"
)

type GinOAuth2Verifier struct {
	GithubConfig *oauth2.Config
	GoogleConfig *oauth2.Config
	RedditConfig *oauth2.Config
}

func (ginOauth2Verifier *GinOAuth2Verifier) AuthVerifier() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader("Authorization")
		if header == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		token := header[len("Bearer "):]
		if token == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		oauth2Token := oauth2.Token{
			AccessToken: token,
			TokenType:   "Bearer",
		}

		if strings.HasPrefix(token, "gho_") {
			// Github token
			// https://docs.github.com/en/rest/apps/oauth-applications?apiVersion=2022-11-28#check-a-token
			basicAuthTransport := github.BasicAuthTransport{
				Username: ginOauth2Verifier.GithubConfig.ClientID,
				Password: ginOauth2Verifier.GithubConfig.ClientSecret,
			}
			githubClient := github.NewClient(basicAuthTransport.Client())

			tokenInfo, _, err := githubClient.Authorizations.Check(ctx, ginOauth2Verifier.GithubConfig.ClientID, token)
			if err != nil {
				println(err.Error())
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			s, _ := json.MarshalIndent(tokenInfo, "", "  ")
			log.Printf("Token info: %s\n", string(s))

		} else if strings.HasPrefix(token, "ya29.") {
			// Google token
			oAuth2Service, err := goauth.NewService(ctx, option.WithTokenSource(ginOauth2Verifier.GoogleConfig.TokenSource(ctx, &oauth2Token)))
			if err != nil {
				ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to create oauth service: %w", err))
				return
			}

			tokenInfo, err := oAuth2Service.Tokeninfo().Do()
			if err != nil {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			s, _ := json.MarshalIndent(tokenInfo, "", "  ")
			log.Printf("Token info: %s\n", string(s))
			ctx.Next()

		} else {
			// JWT
		}
	}
}
