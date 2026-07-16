package main

import (
	"golang-proj/internal/server"
	"golang-proj/internal/service"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := InitiateDatabaseConnection()

	if err != nil {
		log.Fatalf("Error initiating database connection. Error: %s", err.Error())
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	tycoonSvc := service.NewTycoonService(db)
	if tycoonSvc == nil {
		log.Fatal("Error injecting db into service.")
		return
	}

	/**
		to modify this, go to internal/server/routes.go
	**/
	router := server.SetupRoutes(tycoonSvc)

	log.Println("Starting server on :" + port + "...")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("%s", err.Error())
		return
	}
}
