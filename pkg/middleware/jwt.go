package middleware

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/sesaquecruz/go-chat-api/pkg/log"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/jwks"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

type JwtAllClaims struct {
	Issuer    string
	Subject   string
	Audience  []string
	Expiry    int64
	NotBefore int64
	IssuedAt  int64
	ID        string
	Nickname  string
}

type JwtCustomClaims struct {
	Nickname string `json:"https://nickname.com"`
}

func (c *JwtCustomClaims) Validate(ctx context.Context) error {
	return nil
}

var customClaims = func() validator.CustomClaims {
	return &JwtCustomClaims{}
}

func JwtMiddleware(issuer string, audience []string) gin.HandlerFunc {
	logger := log.NewLogger("JwtMiddleware")

	issuerURL, err := url.Parse(issuer)
	if err != nil {
		logger.Fatal(err)
	}

	provider := jwks.NewCachingProvider(issuerURL, 10*time.Minute)

	jwtValidator, err := validator.New(
		provider.KeyFunc,
		validator.RS256,
		issuerURL.String(),
		audience,
		validator.WithCustomClaims(customClaims),
	)
	if err != nil {
		logger.Fatal(err)
	}

	jwtErrorHandler := func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error(err)
	}

	jwtMiddleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(jwtErrorHandler),
	)

	return func(c *gin.Context) {
		unauthorized := true

		var jwtValidHandler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			unauthorized = false
			c.Request = r
			c.Next()
		}

		jwtMiddleware.CheckJWT(jwtValidHandler).ServeHTTP(c.Writer, c.Request)

		if unauthorized {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func JwtClaims(c *gin.Context) (*JwtAllClaims, error) {
	claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("fail to get jwt claims")
	}

	registered := claims.RegisteredClaims
	custom := claims.CustomClaims.(*JwtCustomClaims)

	return &JwtAllClaims{
		Issuer:    registered.Issuer,
		Subject:   registered.Subject,
		Audience:  registered.Audience,
		Expiry:    registered.Expiry,
		NotBefore: registered.NotBefore,
		IssuedAt:  registered.IssuedAt,
		ID:        registered.ID,
		Nickname:  custom.Nickname,
	}, nil
}
