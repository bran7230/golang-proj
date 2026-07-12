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
	dbName := os.Getenv("DOCKER_DB_NAME")
	dbPass := os.Getenv("DOCKER_PASS")
	dbPort := os.Getenv("DOCKER_PORT")

	if dbName == "" || dbPass == "" || dbPort == "" {
		log.Fatal("Missing congifurations. Please check .env.")
		return
	}

	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	log.Print(dsn)

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
