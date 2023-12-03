package v1

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"github.com/Blarc/advent-of-code-bingo/internal/usecase"
	"github.com/Blarc/advent-of-code-bingo/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"google.golang.org/appengine/log"
	"io"
	"net/http"
)

type UserAgentHeaderTransport struct {
	UserAgent string
	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil.
	Transport http.RoundTripper
}

type AuthMiddleware struct {
	u                  usecase.User
	l                  logger.Interface
	TokenEncryptSecret string
	GithubConfig       *oauth2.Config
	GoogleConfig       *oauth2.Config
	RedditConfig       *oauth2.Config
}

func NewAuthMiddleware(
	u usecase.User,
	l logger.Interface,
	tokenEncryptSecret string,
	githubConfig, googleConfig, redditConfig *oauth2.Config,
) *AuthMiddleware {
	return &AuthMiddleware{
		u,
		l,
		tokenEncryptSecret,
		githubConfig,
		googleConfig,
		redditConfig,
	}
}

func (a *AuthMiddleware) Verifier() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userUuid := a.getUserUuidFromHeader(ctx)
		if userUuid == nil {
			errorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := a.u.GetUser(*userUuid)
		if err != nil {
			log.Errorf(ctx, "Failed to get user: %v\n", err)
			errorResponse(ctx, http.StatusUnauthorized, "Unauthorized")
			return
		}

		ctx.Set("user", user)
		ctx.Next()
	}
}

func (a *AuthMiddleware) getUserUuidFromHeader(ctx *gin.Context) *uuid.UUID {
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

	decryptedUuid, decryptionErr := a.Decrypt(encryptedUuid)
	if decryptionErr != nil {
		a.l.Error(decryptionErr, "Failed to decrypt token")
		return nil
	}

	userUuid, err := uuid.Parse(decryptedUuid)
	if err != nil {
		a.l.Error(err, "Failed to create UUID from string")
		return nil
	}

	return &userUuid
}

func (a *AuthMiddleware) Encrypt(plainData string) (string, error) {
	cipherBlock, err := aes.NewCipher([]byte(a.TokenEncryptSecret))
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

// Decrypt decrypts encrypt string with a secret key and returns plain string.
func (a *AuthMiddleware) Decrypt(encodedData string) (string, error) {
	encryptData, err := base64.URLEncoding.DecodeString(encodedData)
	if err != nil {
		return "", err
	}

	cipherBlock, err := aes.NewCipher([]byte(a.TokenEncryptSecret))
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
