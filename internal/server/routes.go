package server

import (
	"golang-proj/internal/service"
	"net/http"
)

// SetupRoutes creates routes with database health checks for readiness probe
// we now use the secretKey for the HMAC key fron the .env. Every endpoint(that uses decode) needs this.
func SetupRoutes(tycoonSvc service.TycoonProcessor, apiKey string, secretKey string) *http.ServeMux {
	mux := http.NewServeMux()

	handleSavesEndpoint := AuthMiddleware(HandleSaves(tycoonSvc, []byte(secretKey)), apiKey)
	mux.Handle("POST /", handleSavesEndpoint)

	return mux
}
