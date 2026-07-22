package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository interface {
	InsertUser(dbQuery string, args ...any) error
	HealthCheck(ctx context.Context) error
}

type Database struct {
	db *sql.DB
}

// DbConfig holds database connection pool configuration
type DbConfig struct {
	MaxOpenConns    int           // max open connections in pool (default 25)
	MaxIdleConns    int           // max idle connections to keep (default 10)
	ConnMaxLifetime time.Duration // max lifetime of a connection (default 5 minutes)
	ConnMaxIdleTime time.Duration // max idle time before closing (default 2 minutes)
}

// DefaultDbConfig returns reasonable production defaults
func DefaultDbConfig() DbConfig {
	return DbConfig{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}
}

func NewServer(db *sql.DB) *Database {
	return &Database{db: db}
}

// ConnectToDatabase initializes a PostgreSQL connection pool with the given config
func ConnectToDatabase(dsn string) (*Database, context.Context, error) {
	return ConnectToDatabaseWithConfig(dsn, DefaultDbConfig())
}

// ConnectToDatabaseWithConfig initializes a PostgreSQL connection pool with custom config
func ConnectToDatabaseWithConfig(dsn string, cfg DbConfig) (*Database, context.Context, error) {
	// init connection with pgxpool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, nil, fmt.Errorf("could not initiate connection pool: %w", err)
	}
	db := stdlib.OpenDBFromPool(pool)

	// ping db with context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, nil, fmt.Errorf("failed to ping database: %w", err)
	}

	slog.Debug("Successfully connected to database.")

	// configure connection pool
	if cfg.MaxOpenConns <= 0 {
		cfg.MaxOpenConns = 25
	}
	if cfg.MaxIdleConns <= 0 {
		cfg.MaxIdleConns = 10
	}
	if cfg.ConnMaxLifetime <= 0 {
		cfg.ConnMaxLifetime = 5 * time.Minute
	}
	if cfg.ConnMaxIdleTime <= 0 {
		cfg.ConnMaxIdleTime = 2 * time.Minute
	}

	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	slog.Debug("Database pool configured", slog.Group(
		"databaseConfiguration",
		slog.Any("maxIdleCons", cfg.MaxIdleConns),
		slog.Any("maxOpenCons", cfg.MaxOpenConns),
		slog.Any("maxConLifetime", cfg.ConnMaxLifetime),
		slog.Any("maxConIdleTime", cfg.ConnMaxIdleTime)),
	)

	return NewServer(db), ctx, nil
}

// Close closes the database connection
func (s *Database) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// HealthCheck pings the database to verify connectivity
func (s *Database) HealthCheck(ctx context.Context) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database not initialized")
	}
	return s.db.PingContext(ctx)
}

// InsertUser executes an insert query with parameterized args (prevents SQL injection)
func (s *Database) InsertUser(dbQuery string, args ...any) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database repository is not initialized")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, dbQuery, args...)
	if err != nil {
		return fmt.Errorf("exec failed: %w", err)
	}
	return nil
}
