package gin_oauth2

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang/glog"
	"golang.org/x/oauth2"
	"log"
	"net/http"
)

const (
	stateKey  = "state"
	sessionID = "gin_oauth2_session"
)

type GinOAuth2 struct {
	Config *oauth2.Config
}

func init() {
	gob.Register(oauth2.Token{})
}

func (ginOauth2 *GinOAuth2) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// Check if token is already in session and set it in context if it is
		session := sessions.Default(ctx)
		existingSession := session.Get(sessionID)
		if token, ok := existingSession.(oauth2.Token); ok {
			ctx.Set("token", token)
			ctx.Next()
			return
		}

		// Check if state is valid - this prevents CSRF attacks
		retrievedState := session.Get(stateKey)
		log.Println("Retrieved state: ", retrievedState)
		if retrievedState != ctx.Query(stateKey) {
			ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s.", retrievedState))
			return
		}

		// Exchange code for token
		log.Println("Exchanging code for token")
		token, err := ginOauth2.Config.Exchange(context.TODO(), ctx.Query("code"))
		if err != nil {
			ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to exchange code for oauth token: %w.", err))
			return
		}
		// TODO: Not sure if this is needed
		ctx.Set("token", *token)

		// Save token in session
		log.Println("Saving session")
		session.Set(sessionID, token)
		if err := session.Save(); err != nil {
			log.Printf("[GIN-oauth2] Failed to save session: %v", err)
			ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to save session: %v.", err))
			return
		}
	}
}

func (ginOauth2 *GinOAuth2) LoginRedirectHandler(ctx *gin.Context) {
	stateValue := randToken()
	session := sessions.Default(ctx)
	session.Set(stateKey, stateValue)
	session.Save()
	ctx.Redirect(http.StatusTemporaryRedirect, ginOauth2.Config.AuthCodeURL(stateValue))
}

func (ginOauth2 *GinOAuth2) LoginHandler(ctx *gin.Context, authProviderName string) {
	stateValue := randToken()
	session := sessions.Default(ctx)
	session.Set(stateKey, stateValue)
	session.Save()
	log.Println(ginOauth2.Config.AuthCodeURL(stateValue))
	ctx.Writer.Write([]byte(`
	<html>
	  <body>
			<a href='` + ginOauth2.Config.AuthCodeURL(stateValue) + `'>
				<button>Login with ` + authProviderName + `!</button>
			</a>
		</body>
	</html>`))
}

func randToken() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		glog.Fatalf("[Gin-OAuth] Failed to read rand: %v", err)
	}
	return base64.StdEncoding.EncodeToString(b)
}
