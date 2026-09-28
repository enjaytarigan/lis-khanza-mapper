package auth

import (
	"crypto/subtle"
	"net/http"
	"strings"
)

// APIKeyMiddleware accepts X-API-Key or Authorization: Bearer <key>.
func APIKeyMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if apiKey == "" {
				http.Error(w, `{"status":false,"message":"webhook not configured"}`, http.StatusServiceUnavailable)
				return
			}
			got := strings.TrimSpace(r.Header.Get("X-API-Key"))
			if got == "" {
				authz := strings.TrimSpace(r.Header.Get("Authorization"))
				if strings.HasPrefix(strings.ToLower(authz), "bearer ") {
					got = strings.TrimSpace(authz[7:])
				}
			}
			if got == "" || subtle.ConstantTimeCompare([]byte(got), []byte(apiKey)) != 1 {
				w.Header().Set("Content-Type", "application/json")
				http.Error(w, `{"status":false,"code":401,"message":"Unauthorized"}`, http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
