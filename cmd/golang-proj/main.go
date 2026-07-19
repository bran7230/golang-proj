package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "golang-proj/internal/repository"
	"golang-proj/internal/server"
	"golang-proj/internal/service"
)

func main() {
	db, err := InitiateDatabaseConnection()
	if err != nil {
		log.Fatalf("Error initiating database connection: %v", err)
	}

	// ensure DB closed on exit; log errors but don't os.Exit from deferred cleanup
	defer func() {
		if db == nil {
			return
		}
		if err := db.Close(); err != nil {
			log.Printf("error closing database connection: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// setup data queue
	tycoonSvc, err := service.NewTycoonService(db, 10000, 10)
	if err != nil {
		log.Fatalf("Error creating tycoon service: %v", err)
		return
	}

	// to modify this, go to internal/server/routes.go
	router := server.SetupRoutes(tycoonSvc)

	// create http.Server with sensible timeouts
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// run server in goroutine so we can listen for shutdown signals
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("Starting server on :%s...", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- err
		} else {
			serverErrors <- nil
		}
	}()

	// catch OS signals for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("Received signal %v. Shutting down...", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("Server Shutdown error: %v", err)
		}
	case err := <-serverErrors:
		if err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}

	// cleanup service resources
	// tycoonSvc.Close() has no error return in existing code, keep same call
	tycoonSvc.Close()

	log.Println("Shutdown complete.")
}
