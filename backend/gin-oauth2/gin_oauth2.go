package gin_oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"time"
)

const (
	stateKey = "state"
	cookieID = "gin_oauth2_cookie"
)

type GinOAuth2 struct {
	Config *oauth2.Config
}

func init() {
	gob.Register(oauth2.Token{})
}

func (ginOauth2 *GinOAuth2) LoginRedirectHandler(ctx *gin.Context) {
	oauthState := generateStateOauthCookie(ctx)
	url := ginOauth2.Config.AuthCodeURL(oauthState)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ginOauth2 *GinOAuth2) CallbackHandler(ctx *gin.Context) {
	token := ginOauth2.GetToken(ctx)
	ctx.JSON(http.StatusOK, token)
}

func (ginOauth2 *GinOAuth2) GetToken(ctx *gin.Context) *oauth2.Token {
	// Read oauthState from Cookie
	oauthState, _ := ctx.Cookie(cookieID)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
		You must always provide a non-empty string and validate that it matches the state
		query parameter on your redirect callback.
	*/
	if ctx.Query(stateKey) != oauthState {
		log.Println("Invalid oauth state.")
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s.", oauthState))
		return nil
	}

	log.Printf("Code: %v", ctx.Query("code"))

	// Exchange code for token
	log.Println("Exchanging code for token")
	token, err := ginOauth2.Config.Exchange(context.TODO(), ctx.Query("code"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to exchange code for oauth token: %w.", err))
		return nil
	}

	return token
}

func generateStateOauthCookie(ctx *gin.Context) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: cookieID, Value: state, Expires: expiration}
	http.SetCookie(ctx.Writer, &cookie)

	return state
}
