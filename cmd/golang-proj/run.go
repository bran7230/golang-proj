package main

import (
	"fmt"
	"golang-proj/internal/repository"
	"os"
)

func InitiateDatabaseConnection() (*repository.Database, error) {

	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	if dsn == "" {
		return nil, fmt.Errorf("database connection string not found in .env")
	}

	db, databaseConError := repository.ConnectToDatabase(dsn)
	if databaseConError != nil {
		return nil, fmt.Errorf("Error initializing database connection: %s", databaseConError)
	}

	return db, nil
}
