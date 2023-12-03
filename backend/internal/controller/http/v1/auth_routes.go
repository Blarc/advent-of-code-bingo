package v1

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"log"
	"math"
	"net/http"
	"time"
)

const (
	StateCookieId = "oauth_state"
	TokenCookieId = "token"
)

type authRoutes struct {
	u usecase.User
	l logger.Interface
	a *AuthMiddleware
}

func newAuthRoutes(handler *gin.RouterGroup, u usecase.User, l logger.Interface, a *AuthMiddleware) {
	r := &authRoutes{u, l, a}

	h := handler.Group("/auth")
	{
		h.GET("/github", r.githubAuth)
		h.GET("/github/callback", r.githubCallback)

		// TODO: Add Google and Reddit auth
		//h.GET("/google", r.googleAuth)
		//h.GET("/google/callback", r.googleCallback)
		//
		//h.GET("/reddit", r.facebookAuth)
		//h.GET("/reddit/callback", r.facebookCallback)

		h.GET("/logout", func(ctx *gin.Context) {
			ctx.SetCookie(TokenCookieId, "", -1, "/", ctx.Request.Host, true, false)
			ctx.Redirect(http.StatusTemporaryRedirect, "/")
		})
	}
}

func (a *authRoutes) githubAuth(ctx *gin.Context) {
	oauthState := generateStateOauthCookie(ctx)
	url := a.a.GithubConfig.AuthCodeURL(oauthState)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (a *authRoutes) githubCallback(ctx *gin.Context) {
	err := checkOAuthState(ctx)
	if err != nil {
		a.l.Error(err, "http - v1 - githubCallback")
		errorResponse(ctx, http.StatusUnauthorized, err.Error())
		return
	}

	// Exchange code for token
	log.Println("Exchanging code for token")
	ctx2 := context.TODO()
	token, err := a.a.GithubConfig.Exchange(ctx2, ctx.Query("code"))
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, "Failed to exchange code for oauth token.")
	}

	client := github.NewClient(a.a.GithubConfig.Client(context.TODO(), token))
	githubUserData, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, "Failed to get user.")
		return
	}

	userUuid, err := a.u.CreateGithubUser(githubUserData)
	if err != nil {
		errorResponse(ctx, http.StatusUnauthorized, "Failed to create user.")
		return
	}

	encryptedUuid, encryptionErr := a.a.Encrypt(userUuid.String())
	if encryptionErr != nil {
		a.l.Error(encryptionErr, "Failed to encrypt token")
		errorResponse(ctx, http.StatusInternalServerError, "Failed to encrypt token.")
		return
	}
	ctx.SetCookie(TokenCookieId, encryptedUuid, math.MaxInt32, "/", ctx.Request.Host, true, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")

}

func checkOAuthState(ctx *gin.Context) error {
	// Read oauthState from Cookie
	oauthState, _ := ctx.Cookie(StateCookieId)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
		You must always provide a non-empty string and validate that it matches the state
		query parameter on your redirect callback.
	*/
	if ctx.Query("state") != oauthState {
		return fmt.Errorf("invalid session state")
	}

	return nil
}

func generateStateOauthCookie(ctx *gin.Context) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: StateCookieId, Value: state, Expires: expiration}
	http.SetCookie(ctx.Writer, &cookie)

	return state
}
