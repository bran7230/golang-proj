package main

import (
	"context"
	"errors"
	"golang-proj/internal/server"
	"golang-proj/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	logEnv := strings.ToLower(os.Getenv("LOGGING_LEVEL"))

	var loggingLevel slog.Level

	switch logEnv {
	case "prod":
		loggingLevel = slog.LevelInfo

	case "dev":
		loggingLevel = slog.LevelDebug

	default:
		loggingLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     loggingLevel,
	}))

	slog.SetDefault(logger)

	if logEnv == "" {
		slog.Warn("Logging level not configured in env. Defaulting to error logs.")
	} else if logEnv != "dev" && logEnv != "prod" {
		slog.Warn("Unrecognized logging level in env. Defaulting to info logs.", "provided_level", logEnv)
	}

	// init from .env vars
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		slog.Warn("Port not found in env. Defaulting port", "portDefault", port)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		slog.Error("Api key not found in .env")
		os.Exit(1)
	}

	secretKey := os.Getenv("HMAC_KEY")
	if secretKey == "" {
		slog.Error("Cannot find the HMAC key")
		os.Exit(1)
	}

	db, ctx, err := InitiateDatabaseConnection()
	if err != nil {
		slog.LogAttrs(ctx, slog.LevelError, "Error initiating database connection", slog.Group(
			"database connection error",
			slog.String("error", err.Error()),
		))
		os.Exit(1)
	}

	// ensure DB closed on exit; log errors but don't os.Exit from deferred cleanup
	defer func() {
		if db == nil {
			return
		}
		if err := db.Close(); err != nil {
			slog.Error("Error closing database connection.", "error", err)
			return
		}
	}()

	// setup data queue
	tycoonSvc, err := service.NewTycoonService(db, 10000, 10)
	if err != nil {
		slog.Error("Error creating tycoon service.", "error", err)
		//manually close db
		err := db.Close()
		if err != nil {
			slog.Error("Error closing database connection.", "error", err)
			os.Exit(1)
		}
		os.Exit(1)
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
		slog.Debug("Starting server.", "port", port)
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
		slog.Info("Received signal. Shutting down...", "signal", sig)
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			slog.Error("Server Shutdown error.", "error", err.Error())
			os.Exit(1)
		}
	case err := <-serverErrors:
		if err != nil {
			slog.Error("Fatal server error", "error", err.Error())
		}
	}

	// cleanup service resources
	// tycoonSvc.Close() has no error return in existing code, keep same call
	tycoonSvc.Close()

	slog.Info("Server shut down complete.")
}
