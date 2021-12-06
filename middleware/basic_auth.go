package middleware

import (
	"net/http"
)

// BasicAuth is a concrete struct of basic auth verifier.
type BasicAuth struct {
	username, password string
}

// NewBasicAuth is a constructor.
func NewBasicAuth(username, password string) RouteMiddleware {
	return &BasicAuth{username, password}
}

// Verify will verify the request to ensure it comes with an authorized basic auth token.
func (ba *BasicAuth) Verify(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !(username == ba.username && password == ba.password) {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	})
}
