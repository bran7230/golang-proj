package server

import (
	"fmt"
	"net/http"
)

func AuthMiddleware(originalHandler http.Handler, apiKey string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != apiKey {
			err := SendError(w, http.StatusUnauthorized, "error, invalid auth key")
			if err != nil {
				fmt.Printf("auth middleware error: %s\n", err.Error())
				return
			}
			return
		}
		originalHandler.ServeHTTP(w, r)
	})
}
