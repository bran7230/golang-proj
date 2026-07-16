package main

import (
	"fmt"
	"golang-proj/internal/repository"
	"golang-proj/internal/server"
	"log"
	"net/http"
	"os"
)

func StartServer() error {

	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	if dsn == "" {
		return fmt.Errorf("database connection string not found in .env")
	}

	// currently I do not use this, so it's a empty variable.
	_, databaseConError := repository.ConnectToDatabase(dsn)
	if databaseConError != nil {
		return fmt.Errorf("Error initializing database connection: %s", databaseConError)
	}

	// pass it to a context pool of connections?? add db to the repository.ConnectToDatabase line

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	/**
		to modify this, go to internal/server/routes.go
	**/
	router := server.SetupRoutes()

	log.Println("Starting server on :" + port + "...")
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}
