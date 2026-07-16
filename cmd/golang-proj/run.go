package main

import (
	"fmt"
	"golang-proj/internal/repository"
	"golang-proj/internal/server"
	"golang-proj/internal/service"
	"log"
	"net/http"
	"os"
)

func StartServer() error {

	dsn := os.Getenv("DATABASE_CONNECTION_STRING")
	if dsn == "" {
		return fmt.Errorf("database connection string not found in .env")
	}

	db, databaseConError := repository.ConnectToDatabase(dsn)
	if databaseConError != nil {
		return fmt.Errorf("Error initializing database connection: %s", databaseConError)
	}

	tycoonSvc := service.NewTycoonService(db)
	if tycoonSvc == nil {
		return fmt.Errorf("Error injecting db into service.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	/**
		to modify this, go to internal/server/routes.go
	**/
	router := server.SetupRoutes(tycoonSvc)

	log.Println("Starting server on :" + port + "...")
	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		return fmt.Errorf("%s", err.Error())
	}

	return nil
}
