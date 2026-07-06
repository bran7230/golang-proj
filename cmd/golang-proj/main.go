package main

import (
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/v1/saveData", handleSaves)

	http.ListenAndServe(":8080", mux)
}
