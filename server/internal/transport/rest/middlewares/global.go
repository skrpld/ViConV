package middlewares

import (
	"net/http"

	"github.com/google/uuid"
)

func GlobalMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = uuid.New().String()
		}
		w.Header().Set("X-Request-ID", requestId)

		next.ServeHTTP(w, r)
	})
}
