package auth

import (
	"crypto/subtle"
	"net/http"
)

type Credentials struct {
	Username string
	Password string
}

func Middleware(creds Credentials) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()
			if !ok ||
				subtle.ConstantTimeCompare([]byte(user), []byte(creds.Username)) != 1 ||
				subtle.ConstantTimeCompare([]byte(pass), []byte(creds.Password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="lis-khanza-mapper"`)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
