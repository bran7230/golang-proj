package server

import (
	"golang-proj/internal/service"
	"net/http"
)

// SetupRoutes SetupRoutesWithDB creates routes with database health checks for readiness probe
func SetupRoutes(tycoonSvc service.TycoonProcessor) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", HandleSaves(tycoonSvc))

	return mux
}
