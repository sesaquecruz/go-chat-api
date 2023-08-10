package test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

type Auth0Server struct {
	signer jose.Signer
	host   string
	server *http.Server
}

func NewAuth0Server() *Auth0Server {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		log.Fatal(err)
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

	return &Auth0Server{
		signer: signer,
		host:   host,
		server: &server,
	}
}

func (s *Auth0Server) GetIssuer() string {
	return s.host
}

func (s *Auth0Server) GetAudience() string {
	return s.host + "/userinfo"
}

func (s *Auth0Server) GenerateSub() string {
	return fmt.Sprintf("auth0|%s", strings.ReplaceAll(uuid.NewString(), "-", "")[:24])
}

func (s *Auth0Server) GenerateJWT(subject string) (string, error) {
	claims := jwt.Claims{
		Issuer:   s.GetIssuer(),
		Audience: []string{s.GetAudience()},
		Subject:  subject,
	}

	return jwt.Signed(s.signer).Claims(claims).CompactSerialize()
}

func (s *Auth0Server) Run() error {
	return s.server.ListenAndServe()
}

func (s *Auth0Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
