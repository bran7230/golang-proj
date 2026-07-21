package main

import (
	"context"
	"errors"
	"log/slog"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	}))

	slog.SetDefault(logger)
	// init from .env vars
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	apiKey := os.Getenv("APIKEY")
	if apiKey == "" {
		slog.Error("api key not found in .env")
		return
	}

	secretKey := os.Getenv("HMAC_KEY")
	if secretKey == "" {
		log.Fatalf("Cannot find the HMAC key")
		return
	}

	db, err := InitiateDatabaseConnection()
	if err != nil {
		slog.Error("Error initiating database connection", slog.Group(
			"database connection error",
			slog.String("error", err.Error()),
		))
		return
	}


	// ensure DB closed on exit; log errors but don't os.Exit from deferred cleanup
	defer func() {
		if db == nil {
			return
		}
		if err := db.Close(); err != nil {
			slog.Error("error closing database connection.", "error", err.Error())
			return
		}
	}()

	// setup data queue
	tycoonSvc, err := service.NewTycoonService(db, 10000, 10)
	if err != nil {
		slog.Error("Error creating tycoon service.", "error", err.Error())
		return
	}

	// to modify this, go to internal/server/routes.go
	router := server.SetupRoutes(tycoonSvc, apiKey, secretKey)

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
		slog.Info("Starting server.", "port", port)
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
		slog.Debug("Received signal. Shutting down...", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server Shutdown error.", "error", err.Error())
			return
		}
	case err := <-serverErrors:
		if err != nil {
			slog.Error("Fatal server error", "error", err.Error())
		}
	}

	// cleanup service resources
	// tycoonSvc.Close() has no error return in existing code, keep same call
	tycoonSvc.Close()

	slog.Debug("Server shut down complete.")
}
