package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository interface {
	InsertUser(dbQuery string, args ...any) error
}

type Database struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Database {
	return &Database{db: db}
}

func ConnectToDatabase(dsn string) (*Database, error) {
	// init connection with connection pool
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("could not initiate connection pool. Error: %s", err.Error())
	}
	db := stdlib.OpenDBFromPool(pool)

	// ping db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to database.")

	// limit connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	// set a reasonable max lifetime for connections
	db.SetConnMaxLifetime(5 * time.Minute)

	return NewServer(db), nil
}

// Close expose the connection closure to the rest of the application
func (s *Database) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Database) InsertUser(dbQuery string, args ...any) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database repository is not initialized")
	}
	_, err := s.db.Exec(dbQuery, args...)
	if err != nil {
		return err
	}
	return nil
}
