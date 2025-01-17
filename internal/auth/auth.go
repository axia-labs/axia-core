package auth

import (
	"context"
	"crypto/subtle"
	"errors"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	ErrMissingSecretKey = errors.New("AXIA_SECRET_KEY environment variable is not set")
	ErrInvalidSecretKey = errors.New("invalid secret key")
)

// Authenticator handles API authentication
type Authenticator struct {
	secretKey string
	logger    *logrus.Logger
}

func NewAuthenticator(logger *logrus.Logger) (*Authenticator, error) {
	secretKey := os.Getenv("AXIA_SECRET_KEY")
	if secretKey == "" {
		return nil, ErrMissingSecretKey
	}

	return &Authenticator{
		secretKey: secretKey,
		logger:    logger,
	}, nil
}

// ValidateKey checks if the provided key matches the secret key
func (a *Authenticator) ValidateKey(key string) bool {
	return subtle.ConstantTimeCompare([]byte(key), []byte(a.secretKey)) == 1
}

// Middleware provides HTTP authentication middleware
func (a *Authenticator) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-Axia-Key")
		if key == "" {
			a.logger.Warn("Request missing authentication key")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !a.ValidateKey(key) {
			a.logger.Warn("Invalid authentication key provided")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// ValidateContext validates authentication for non-HTTP operations
func (a *Authenticator) ValidateContext(ctx context.Context) error {
	key, ok := ctx.Value("secret_key").(string)
	if !ok || key == "" {
		return ErrMissingSecretKey
	}

	if !a.ValidateKey(key) {
		return ErrInvalidSecretKey
	}

	return nil
} 