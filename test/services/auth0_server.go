package services

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/sesaquecruz/go-chat-api/pkg/log"

	"github.com/google/uuid"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

type Claims struct {
	jwt.Claims
	Issuer    string           `json:"iss,omitempty"`
	Subject   string           `json:"sub,omitempty"`
	Audience  jwt.Audience     `json:"aud,omitempty"`
	Expiry    *jwt.NumericDate `json:"exp,omitempty"`
	NotBefore *jwt.NumericDate `json:"nbf,omitempty"`
	IssuedAt  *jwt.NumericDate `json:"iat,omitempty"`
	ID        string           `json:"jti,omitempty"`
	Nickname  string           `json:"https://nickname.com"`
}

type Auth0Server struct {
	signer jose.Signer
	host   string
	server *http.Server
	logger *log.Logger
}

func NewAuth0Server() *Auth0Server {
	logger := log.NewLogger("Auth0Server")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Fatal(err)
	}

	algorithm := jose.RS256

	webKey := jose.JSONWebKey{
		Key:       key,
		KeyID:     "kid",
		Algorithm: string(algorithm),
		Use:       "sig",
	}

	signKey := jose.SigningKey{
		Key:       webKey,
		Algorithm: algorithm,
	}

	signer, err := jose.NewSigner(signKey, (&jose.SignerOptions{}).WithType("JWT"))
	if err != nil {
		logger.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		logger.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		logger.Fatal(err)
	}

	host := fmt.Sprintf("http://127.0.0.1:%d", port)
	oidcUri := "/.well-known/openid-configuration"
	jwksUri := "/.well-known/jwks.json"

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case oidcUri:
			res := struct {
				JwksUri string `json:"jwks_uri"`
			}{
				JwksUri: host + jwksUri,
			}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		case jwksUri:
			res := jose.JSONWebKeySet{
				Keys: []jose.JSONWebKey{webKey.Public()},
			}
			if err := json.NewEncoder(w).Encode(res); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		default:
			http.Error(w, "invalid uri", http.StatusNotFound)
		}
	})

	server := http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		server.ListenAndServe()
	}()

	return &Auth0Server{
		signer: signer,
		host:   host,
		server: &server,
		logger: logger,
	}
}

func (s *Auth0Server) GetIssuer() string {
	return s.host
}

func (s *Auth0Server) GetAudience() string {
	return s.host + "/userinfo"
}

func (s *Auth0Server) GetNickname() string {
	return "An username"
}

func (s *Auth0Server) GenerateSub() string {
	return fmt.Sprintf("auth0|%s", strings.ReplaceAll(uuid.NewString(), "-", "")[:24])
}

func (s *Auth0Server) GenerateJWT(subject string) (string, error) {
	claims := Claims{
		Issuer:   s.GetIssuer(),
		Audience: []string{s.GetAudience()},
		Subject:  subject,
		Nickname: s.GetNickname(),
	}

	token, err := jwt.Signed(s.signer).Claims(claims).CompactSerialize()
	if err != nil {
		s.logger.Error(err)
	}

	return token, err
}

func (s *Auth0Server) Terminate(ctx context.Context) error {
	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Error(err)
	}

	return err
}
