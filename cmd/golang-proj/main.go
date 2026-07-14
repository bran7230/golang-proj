package main

import (
	"log"
	"net/http"
	"os"

	// routes, handlers
	"golang-proj/internal/server"

	// db connections
	"golang-proj/internal/repository"
)

func main() {
	/**
		You can move this to a new file, ie: a dedicated run file, and change the .env names to be your own if needed.
	**/
	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	if dsn == "" {
		log.Fatal("Missing database connection string in .env")
	}

	// currently I do not use this, so it's a empty variable.
	_, databaseConError := repository.ConnectToDatabase(dsn)
	if databaseConError != nil {
		log.Fatal("Error initializing database connection: ", databaseConError)
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}
	/**
		to modify this, go to internal/server
	**/
	router := server.SetupRoutes()

	log.Println("Starting server on :" + port + "...")
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		log.Fatal("Error during server setup ..." + err.Error())
	}

}
