package middleware

import (
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

func JwtClaims(c *gin.Context) (*validator.RegisteredClaims, error) {
	claims, ok := c.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
	if !ok {
		return nil, errors.New("fail to get jwt claims")
	}

	return &claims.RegisteredClaims, nil
}
