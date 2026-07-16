package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Repository interface {
	InsertUser(dbQuery string, args ...interface{}) error
}

type Server struct {
	db *sql.DB
}

func NewServer(db *sql.DB) *Server {
	return &Server{db: db}
}

func ConnectToDatabase(dsn string) (*Server, error) {
	// init connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

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

func (s *Server) InsertUser(dbQuery string, args ...interface{}) error {
	if s == nil || s.db == nil {
		return fmt.Errorf("database repository is not initialized")
	}

	_, err := s.db.Exec(dbQuery, args...)
	if err != nil {
		return err
	}
	return nil
}
