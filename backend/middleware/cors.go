package middleware

import (
	"net/http"
	"strings"
)

// CORS returns a middleware wrapper with a configurable allowed origin for production security.
// allowedOrigin may be a single origin or a comma-separated list. The request's Origin is
// echoed back when it is in the allowlist. Requests from origins outside the allowlist never
// receive an Access-Control-Allow-Origin header, so browsers correctly reject them.
func CORS(allowedOrigin string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed := parseOrigins(allowedOrigin)
			origin := r.Header.Get("Origin")
			w.Header().Set("Vary", "Origin")

			if origin != "" && !containsOrigin(allowed, origin) {
				// Do not echo a different allowed origin for a disallowed request.
				// Reject preflights so the browser reports the real CORS failure.
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
				return
			}

			if origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// parseOrigins splits a comma-separated origin allowlist, falling back to the
// local Vite dev servers used by this project.
func parseOrigins(value string) []string {
	var origins []string
	for _, part := range strings.Split(value, ",") {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000", "http://localhost:5173"}
	}
	return origins
}

func containsOrigin(origins []string, origin string) bool {
	if origin == "" {
		return false
	}
	for _, allowed := range origins {
		if strings.EqualFold(strings.TrimRight(allowed, "/"), strings.TrimRight(origin, "/")) {
			return true
		}
	}
	return false
}
