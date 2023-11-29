package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/Blarc/advent-of-code-bingo/models"
	"github.com/Blarc/advent-of-code-bingo/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/v56/github"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	goauth "google.golang.org/api/oauth2/v2"
	"google.golang.org/api/option"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	stateCookieId = "oauth_state"
	tokenCookieId = "token"
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

func (o *OAuth) GetToken(ctx *gin.Context) *oauth2.Token {
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
		return nil
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
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Failed to exchange code for oauth token: %w.", err))
		return nil
	}
	return token
}

func GithubCallbackHandler(ctx *gin.Context, config *OAuth) {
	token := config.GetToken(ctx)

	client := github.NewClient(config.Config.Client(context.TODO(), token))
	githubUserData, _, err := client.Users.Get(context.TODO(), "")
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Failed to get user: %v", err))
		return
	}
	githubId := strconv.FormatInt(*githubUserData.ID, 10)

	var user models.User
	result := models.DB.Where(models.User{GithubID: githubId}).Assign(models.User{
		GithubID:  githubId,
		Name:      *githubUserData.Name,
		AvatarURL: *githubUserData.AvatarURL,
		GithubURL: *githubUserData.HTMLURL,
	}).FirstOrCreate(&user)

	if result.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to create user: %v", result.Error))
		return
	}

	encryptedUuid, err := encrypt(user.ID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to encrypt user id: %v", err))
		return
	}
	ctx.SetCookie(tokenCookieId, encryptedUuid, math.MaxInt32, "/", ctx.Request.Host, true, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func GoogleCallbackHandler(ctx *gin.Context, config *OAuth) {
	token := config.GetToken(ctx)

	service, err := goauth.NewService(ctx, option.WithTokenSource(config.Config.TokenSource(ctx, token)))
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	googleUserData, err := service.Userinfo.Get().Do()
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Failed to get user: %v", err))
		return
	}

	var user models.User
	result := models.DB.Where(models.User{GoogleID: googleUserData.Id}).Assign(models.User{
		GoogleID:  googleUserData.Id,
		Name:      googleUserData.Name,
		AvatarURL: googleUserData.Picture,
	}).FirstOrCreate(&user)

	if result.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to create user: %v", result.Error))
		return
	}

	encryptedUuid, err := encrypt(user.ID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to encrypt user id: %v", err))
		return
	}
	ctx.SetCookie(tokenCookieId, encryptedUuid, math.MaxInt32, "/", ctx.Request.Host, true, false)
	ctx.Redirect(http.StatusTemporaryRedirect, "/")
}

func RedditCallbackHandler(ctx *gin.Context, config *OAuth) {
	token := config.GetToken(ctx)

	log.Println(token.AccessToken)

	req, err := http.NewRequest("GET", "https://oauth.reddit.com/api/v1/me", nil)
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to create request: %v", err))
		return
	}

	req.Header.Set("Authorization", "Bearer "+token.AccessToken)
	req.Header.Set("User-Agent", config.UserAgent)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		ctx.AbortWithError(http.StatusUnauthorized, fmt.Errorf("Failed to get user: %v", err))
		return
	}

	type RedditSubredditData struct {
		Url string `json:"url"`
	}

	type RedditUserData struct {
		ID        string              `json:"id"`
		Name      string              `json:"name"`
		IconImg   string              `json:"icon_img"`
		Subreddit RedditSubredditData `json:"subreddit"`
	}

	redditUserDataRaw, err := io.ReadAll(response.Body)
	var redditUserData RedditUserData
	json.Unmarshal(redditUserDataRaw, &redditUserData)

	var user models.User
	result := models.DB.Where(models.User{RedditID: redditUserData.ID}).Assign(models.User{
		RedditID:  redditUserData.ID,
		Name:      redditUserData.Name,
		AvatarURL: redditUserData.IconImg,
		RedditURL: "https://www.reddit.com" + redditUserData.Subreddit.Url,
	}).FirstOrCreate(&user)

	if result.Error != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to create user: %v", result.Error))
		return
	}

	encryptedUuid, err := encrypt(user.ID.String())
	if err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, fmt.Errorf("Failed to encrypt user id: %v", err))
		return
	}
	ctx.SetCookie(tokenCookieId, encryptedUuid, math.MaxInt32, "/", ctx.Request.Host, true, false)
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

func Verifier() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userUuid := GetUserUuidFromHeader(ctx)
		if userUuid == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Printf("Searching for user with id %s\n", userUuid)
		var user models.User
		result := models.DB.
			Preload("BingoCards").
			Preload("BingoBoards").
			Preload("PersonalBingoBoard").
			First(&user, "id = ?", userUuid)
		if result.Error != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		log.Printf("User successfully logged in: %+v\n", user)
		ctx.Set("user", user)
		ctx.Next()
	}
}

func GetUserUuidFromHeader(ctx *gin.Context) *uuid.UUID {
	header := ctx.GetHeader("Authorization")
	if header == "" {
		return nil
	}

	if len(header) <= len("Bearer ") {
		return nil
	}

	encryptedUuid := header[len("Bearer "):]
	if encryptedUuid == "" {
		return nil
	}

	decryptedUuid, decryptionErr := decrypt(encryptedUuid)
	if decryptionErr != nil {
		log.Printf("Failed to decrypt token: %v\n", decryptionErr)
		return nil
	}

	userUuid, err := uuid.Parse(decryptedUuid)
	if err != nil {
		log.Printf("Failed to create UUID from string: %v\n", err)
		return nil
	}

	return &userUuid
}

func encrypt(plainData string) (string, error) {
	cipherBlock, err := aes.NewCipher([]byte(utils.GetEnvVariable("TOKEN_ENCRYPT_SECRET")))
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(plainData), nil)), nil
}

// decrypt decrypts encrypt string with a secret key and returns plain string.
func decrypt(encodedData string) (string, error) {
	encryptData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher([]byte(utils.GetEnvVariable("TOKEN_ENCRYPT_SECRET")))
	if err != nil {
		return "", err
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		return "", err
	}

	nonceSize := aead.NonceSize()
	if len(encryptData) < nonceSize {
		return "", err
	}

	nonce, cipherText := encryptData[:nonceSize], encryptData[nonceSize:]
	plainData, err := aead.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainData), nil
}
