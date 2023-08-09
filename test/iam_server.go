package test

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/sesaquecruz/go-chat-api/pkg"

	"github.com/google/uuid"
	"gopkg.in/go-jose/go-jose.v2"
	"gopkg.in/go-jose/go-jose.v2/jwt"
)

type IamServer struct {
	webKey   jose.JSONWebKey
	signer   jose.Signer
	issuer   string
	audience string
	server   *http.Server
	running  bool
	logger   *pkg.Logger
}

func NewIamServer() *IamServer {
	logger := pkg.NewLogger("IamServer")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		logger.Error(err)
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
		logger.Error(err)
	}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		logger.Error(err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		logger.Error(err)
	}

	issuer := fmt.Sprintf("http://127.0.0.1:%d", port)

	oidcUri := "/.well-known/openid-configuration"
	jwksUri := "/.well-known/jwks.json"

	var handler http.Handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case oidcUri:
			res := struct {
				JwksUri string `json:"jwks_uri"`
			}{
				JwksUri: issuer + jwksUri,
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

	return &IamServer{
		webKey:   webKey,
		signer:   signer,
		issuer:   issuer,
		audience: issuer + "/userinfo",
		server:   &server,
		logger:   logger,
	}
}

func (s *IamServer) GetIssuer() string {
	return s.issuer
}

func (s *IamServer) GetAudience() string {
	return s.audience
}

func (s *IamServer) GenerateSub() string {
	return fmt.Sprintf("auth0|%s", strings.ReplaceAll(uuid.NewString(), "-", "")[:24])
}

func (s *IamServer) GenerateJWT(subject string) (string, error) {
	claims := jwt.Claims{
		Issuer:   s.issuer,
		Audience: []string{s.audience},
		Subject:  subject,
	}

	return jwt.Signed(s.signer).Claims(claims).CompactSerialize()
}

func (s *IamServer) IsRunning() bool {
	return s.running
}

func (s *IamServer) Run() error {
	if !s.running {
		s.running = true
		s.logger.Info("running iam server")
		return s.server.ListenAndServe()
	}

	return errors.New("iam server already running")
}

func (s *IamServer) Stop(ctx context.Context) error {
	if s.running {
		s.running = false
		s.logger.Info("stopping iam server")
		return s.server.Shutdown(ctx)
	}

	return errors.New("aim server is not running")
}
