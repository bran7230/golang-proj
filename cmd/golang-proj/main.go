package main

import (
	"golang-proj/internal/repository"
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

	defer func(db *repository.Database) {
		err := db.Close()
		if err != nil {
			log.Fatalf("error closing database connection. Error: %s", err.Error())
		}
	}(db)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// setup data queue
	tycoonSvc := service.NewTycoonService(db, 10000, 10)
	if tycoonSvc == nil {
		log.Fatal("Error injecting db into service.")
		return
	}

	defer tycoonSvc.Close()

	// to modify this, go to internal/server/routes.go
	router := server.SetupRoutes(tycoonSvc)

	log.Println("Starting server on :" + port + "...")

	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("%s", err.Error())
		return
	}
}
