package middleware

import (
	"encoding/base64"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"micro/pkg/logger"
	"micro/pkg/util"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	// len of authorization header
	lenOfAuthorization = 2

	// AuthTypeBasic specify authentication process with Basic auth
	AuthTypeBasic = "Basic"

	// AuthTypeBasicAppKey specify authentication process with Basic auth + Application Key
	AuthTypeBasicAppKey = "BasicKey"

	// AuthTypeBearer specify authentication process with Bearer token
	AuthTypeBearer = "Bearer"
)

// Guard handle authentication and authorization process.
type Guard struct {
	logger *logger.Logger

	authenticationMethod []string
}

// NewGuard will initialize a new Guard middleware.
func NewGuard(logger *logger.Logger) *Guard {
	return &Guard{
		logger:               logger,
		authenticationMethod: []string{},
	}
}

// SetAuthenticationMethod will set default authentication method.
func (g *Guard) SetAuthenticationMethod(methods ...string) {
	g.authenticationMethod = append(g.authenticationMethod, methods...)
}

// Authenticate uses to handle authentication process.
func (g *Guard) Authenticate(opts ...GuardOption) gin.HandlerFunc {
	for _, opt := range opts {
		opt(g)
	}

	return func(c *gin.Context) {
		var authType string

		headerAppKey := c.Request.Header.Get("Application-Key")
		headerAuth := c.Request.Header.Get("Authorization")
		headerAuthType := strings.SplitN(headerAuth, " ", 2)

		if len(headerAuthType) != lenOfAuthorization {
			g.logger.Log.Infof("Error validating given authorization header")
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.unauthorized"))
			return
		}

		if len(headerAuthType) == lenOfAuthorization {
			authType = headerAuthType[0]
		}

		if len(g.authenticationMethod) == 0 {
			g.logger.Log.Infof("Error validating given authorization type, err: no authentication method are specified")
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.unauthorized"))
			return
		}

		if authType == AuthTypeBasic && headerAppKey != "" {
			authType = AuthTypeBasicAppKey
		}

		if !util.SliceContains(g.authenticationMethod, authType) {
			g.logger.Log.Infof("Error validating given authorization type, err: auth type %s is not allowed, allowed method: %s", authType, strings.Join(g.authenticationMethod, ", "))
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.unauthorized"))
			return
		}

		if authType == AuthTypeBasic {
			payload, _ := base64.StdEncoding.DecodeString(headerAuthType[1])
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != lenOfAuthorization {
				g.logger.Log.Infof("Error validating given authorization header, err: invalid basic auth")
				_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.auth.invalid_basic_auth"))
				return
			}

			c.Set("auth_type", AuthTypeBasic)
		}

		if authType == AuthTypeBasicAppKey {
			payload, _ := base64.StdEncoding.DecodeString(headerAuthType[1])
			pair := strings.SplitN(string(payload), ":", 2)
			if len(pair) != lenOfAuthorization {
				g.logger.Log.Infof("Error validating given authorization header, err: invalid basic auth")
				_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.auth.invalid_basic_auth"))
				return
			}

			if headerAppKey == "" {
				g.logger.Log.Infof("Error validating given authorization header, err: application key is not provided")
				_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.auth.invalid_basic_auth"))
				return
			}

			c.Set("auth_type", AuthTypeBasicAppKey)
		}

		if authType == AuthTypeBearer {
			c.Set("auth_type", AuthTypeBearer)
		}

		return
	}
}

// WithAccessToken uses to set additional Access-Token validation retrieved from the request header.
func (g *Guard) WithAccessToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		headerAccessToken := c.Request.Header.Get("Access-Token")
		if headerAccessToken == "" {
			g.logger.Log.Infof("Error validating authorization header, err: access-token is not provided")
			_ = c.AbortWithError(http.StatusUnauthorized, errors.New("common.error.unauthorized"))
			return
		}

		return
	}
}

// Authorize uses to handle authorization process.
func (g *Guard) Authorize() gin.HandlerFunc {
	return func(c *gin.Context) {
		return
	}
}

// validateJWTToken will validate the given token.
func (g *Guard) validateJWTToken(token string) (interface{}, error) {
	publicKey, _ := base64.StdEncoding.DecodeString(os.Getenv("OAUTH2_ACCESS_TOKEN_PUBLIC_KEY"))
	key, err := jwt.ParseRSAPublicKeyFromPEM(publicKey)
	if err != nil {
		g.logger.Log.Infof("Error parsing RSA public key, err: %v", err)
		return "", errors.New("common.error.unauthorized")
	}

	tok, err := jwt.Parse(token, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodRSA); !ok {
			g.logger.Log.Infof("Error parsing JWT token, unexpected signing method: %s", jwtToken.Header["alg"])
			return nil, errors.New("common.error.unauthorized")
		}

		return key, nil
	})
	if err != nil {
		g.logger.Log.Infof("Error parsing JWT token, err: %v", err)
		return nil, errors.New("common.error.unauthorized")
	}

	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok || !tok.Valid {
		g.logger.Log.Infof("Error claiming JWT claims")
		return nil, errors.New("common.error.unauthorized")
	}

	return claims, nil
}
