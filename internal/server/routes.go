package server

import (
	"context"
	"net/http"
	"time"

	"golang-proj/internal/repository"
	"golang-proj/internal/service"
)

// SetupRoutesWithDB creates routes with database health checks for readiness probe
func SetupRoutesWithDB(tycoonSvc service.TycoonProcessor, db *repository.Database) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", HandleSaves(tycoonSvc))

	// liveness probe
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// readiness probe with database health check
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, r *http.Request) {
		if tycoonSvc == nil {
			http.Error(w, "service not ready", http.StatusServiceUnavailable)
			return
		}

		if db == nil {
			http.Error(w, "database not available", http.StatusServiceUnavailable)
			return
		}

		// check database connectivity with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		if err := db.HealthCheck(ctx); err != nil {
			http.Error(w, "database not healthy", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready"))
	})

	return mux
}
