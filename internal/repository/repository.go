package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository interface {
	InsertUser(dbQuery string, args ...interface{}) error
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
		return nil, fmt.Errorf("Could not initiate connection pool. Error: %s", err.Error())
	}
	db := stdlib.OpenDBFromPool(pool)

	// ping db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Sucessfully connected to database.")

	// limit connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return NewServer(db), nil
}

func (s *Database) InsertUser(dbQuery string, args ...interface{}) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database repository is not initialized")
	}

	_, err := s.db.Exec(dbQuery, args...)
	if err != nil {
		return err
	}
	return nil
}
