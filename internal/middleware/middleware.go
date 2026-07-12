package middleware

import (
	"log"
	"net/http"
	"time"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// call next handler in chain
		next.ServeHTTP(w, r)

		// log the duration
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}
