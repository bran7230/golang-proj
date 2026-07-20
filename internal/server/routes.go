package server

import (
	"golang-proj/internal/service"
	"net/http"
)

// SetupRoutes SetupRoutesWithDB creates routes with database health checks for readiness probe
func SetupRoutes(tycoonSvc service.TycoonProcessor, apiKey string) *http.ServeMux {
	mux := http.NewServeMux()

	handleSavesEndpoint := AuthMiddleware(HandleSaves(tycoonSvc), apiKey)
	mux.Handle("POST /", handleSavesEndpoint)

	return mux
}
