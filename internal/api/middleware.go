package api

import (
	"net/http"
	"os"

	"github.com/mtbuzato/go-challenge/internal/errors"
)

func (s *apiServer) mdwHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func (s *apiServer) mdwAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		key := os.Getenv("API_KEY")

		if auth != "Bearer "+key {
			s.handleError(w, errors.NewHTTPError("You don't have permission to access this endpoint.", http.StatusUnauthorized))
			return
		}

		next.ServeHTTP(w, r)
	})
}
