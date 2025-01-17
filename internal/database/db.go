package database

import (
	"context"
	"fmt"
	"time"
	
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/sirupsen/logrus"
	"axia/internal/axiom"
	"axia/internal/auth"
)

type DB struct {
	pool   *pgxpool.Pool
	logger *logrus.Logger
	auth   *auth.Authenticator
}

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

func New(cfg Config, logger *logrus.Logger) (*DB, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	
	pool, err := pgxpool.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	return &DB{
		pool:   pool,
		logger: logger,
	}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) StoreClaim(ctx context.Context, claim *axiom.Claim) error {
	if err := db.auth.ValidateContext(ctx); err != nil {
		return fmt.Errorf("authentication failed: %w", err)
	}
	// ... rest of the method
} 