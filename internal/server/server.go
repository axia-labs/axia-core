package server

import (
	"context"
	"fmt"
	"net/http"
	"github.com/sirupsen/logrus"
	"axia/internal/axiom"
	"axia/internal/trust"
	"axia/internal/social/twitter"
	"axia/internal/auth"
)

type Server struct {
	server  *http.Server
	logger  *logrus.Logger
	twitter *twitter.Handler
	auth    *auth.Authenticator
}

func NewServer(port int, manager *axiom.Manager, network *trust.Network, logger *logrus.Logger) (*Server, error) {
	auth, err := auth.NewAuthenticator(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize authenticator: %w", err)
	}

	twitterHandler := twitter.NewHandler(manager, network, logger)
	
	mux := http.NewServeMux()
	
	// Wrap handlers with authentication middleware
	mux.Handle("/webhook/twitter", auth.Middleware(http.HandlerFunc(twitterHandler.HandleWebhook)))

	return &Server{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
		logger:  logger,
		twitter: twitterHandler,
		auth:    auth,
	}, nil
}

func (s *Server) Start() error {
	s.logger.WithField("addr", s.server.Addr).Info("Starting HTTP server")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.logger.Info("Shutting down HTTP server")
	return s.server.Shutdown(ctx)
} 