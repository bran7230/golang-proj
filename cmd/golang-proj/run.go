package main

import (
	"context"
	"fmt"
	"golang-proj/internal/repository"
	"os"
)

func InitiateDatabaseConnection() (*repository.Database, context.Context, error) {

	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	if dsn == "" {
		return nil, nil, fmt.Errorf("database connection string not found in .env")
	}

	db, ctx, databaseConError := repository.ConnectToDatabase(dsn)
	if databaseConError != nil {
		return nil, ctx, fmt.Errorf("error initializing database connection: %s", databaseConError)
	}
	return db, ctx, nil
}
