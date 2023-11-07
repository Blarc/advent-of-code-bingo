package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"golang.org/x/oauth2"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"log"
	"math"
	"net/http"
	"strings"
	"time"
)

const (
	stateCookieId = "oauth_state"
	tokenCookieId = "oauth_token"
)

type OAuth struct {
	Config    *oauth2.Config
	UserAgent string
}

func init() {
	gob.Register(oauth2.Token{})
}

func (o *OAuth) LoginRedirectHandler(ctx *gin.Context) {
	oauthState := generateStateOauthCookie(ctx)
	url := o.Config.AuthCodeURL(oauthState)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (o *OAuth) CallbackHandler(ctx *gin.Context) {
	// Read oauthState from Cookie
	oauthState, _ := ctx.Cookie(stateCookieId)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks.
		You must always provide a non-empty string and validate that it matches the state
		query parameter on your redirect callback.
	*/
	if ctx.Query("state") != oauthState {
		log.Println("Invalid oauth state.")
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Invalid session state: %s.", oauthState))
		return
	}

	ctx2 := context.TODO()
	if o.UserAgent != "" {
		transport := UserAgentHeaderTransport{
			UserAgent: o.UserAgent,
		}
		ctx2 = context.WithValue(ctx, oauth2.HTTPClient, transport.Client())
	}

	// Exchange code for token
	log.Println("Exchanging code for token")
	token, err := o.Config.Exchange(ctx2, ctx.Query("code"))
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, fmt.Errorf("Failed to exchange code for oauth token: %w.", err))
		return
	}

	age := int(token.Expiry.Sub(time.Now()).Minutes())
	if age < 0 {
		age = math.MaxInt32
	}

	ctx.SetCookie(tokenCookieId, token.AccessToken, age, "/", ctx.Request.Host, true, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func generateStateOauthCookie(ctx *gin.Context) string {
	var expiration = time.Now().Add(365 * 24 * time.Hour)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: stateCookieId, Value: state, Expires: expiration}
	http.SetCookie(ctx.Writer, &cookie)

	return state
}

type UserAgentHeaderTransport struct {
	UserAgent string
	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

func (t *UserAgentHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := setUserAgentHeader(req, t.UserAgent)
	return t.transport().RoundTrip(req2)
}

// Client returns an *http.Client that makes requests that are authenticated
// using HTTP Basic Authentication.
func (t *UserAgentHeaderTransport) Client() *http.Client {
	return &http.Client{Transport: t}
}

func (t *UserAgentHeaderTransport) transport() http.RoundTripper {
	if t.Transport != nil {
		return t.Transport
	}
	return http.DefaultTransport
}

func setUserAgentHeader(req *http.Request, userAgent string) *http.Request {
	// To set extra headers, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	//
	// Since we are going to modify only req.Header here, we only need a deep copy
	// of req.Header.
	convertedRequest := new(http.Request)
	*convertedRequest = *req
	convertedRequest.Header = make(http.Header, len(req.Header))

	for k, s := range req.Header {
		convertedRequest.Header[k] = append([]string(nil), s...)
	}
	convertedRequest.Header.Set("user-agent", userAgent)
	return convertedRequest
}

type Verifier struct {
	GithubConfig *oauth2.Config
	GoogleConfig *oauth2.Config
	RedditConfig *oauth2.Config
}

func (v *Verifier) AuthVerifier() gin.HandlerFunc {
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
				Username: v.GithubConfig.ClientID,
				Password: v.GithubConfig.ClientSecret,
			}
			githubClient := github.NewClient(basicAuthTransport.Client())

			tokenInfo, _, err := githubClient.Authorizations.Check(ctx, v.GithubConfig.ClientID, token)
			if err != nil {
				println(err.Error())
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			s, _ := json.MarshalIndent(tokenInfo, "", "  ")
			log.Printf("Token info: %s\n", string(s))

		} else if strings.HasPrefix(token, "ya29.") {
			// Google token
			oAuth2Service, err := goauth.NewService(ctx, option.WithTokenSource(v.GoogleConfig.TokenSource(ctx, &oauth2Token)))
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
