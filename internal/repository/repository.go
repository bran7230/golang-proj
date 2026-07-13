package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Server struct {
	db *sql.DB
}

func ConnectToDatabase(dsn string) (*Server, error) {
	// init connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	defer db.Close()

	// ping db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Sucessfully connected to database.")

	// limit connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)

	return &Server{db: db}, nil

}
