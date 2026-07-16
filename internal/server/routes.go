package server

import (
	"net/http"

	"golang-proj/internal/service"
)

func SetupRoutes(tycoonSvc service.TycoonProcessor) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", HandleSaves(tycoonSvc))

	return mux
}
