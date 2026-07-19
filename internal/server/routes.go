package server

import (
	"net/http"

	"golang-proj/internal/service"
)

func SetupRoutes(tycoonSvc service.TycoonProcessor) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", HandleSaves(tycoonSvc))

	// liveness and readiness endpoints
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return mux
}
