package http

import (
	"Tournament/internal/ports/input"
	"context"
	"crypto/subtle"
	"net/http"
)

type userKey struct{}

func AuthMiddleware(userService input.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			username, password, ok := r.BasicAuth()

			if !ok {
				basicAuthFailed(w)
				return
			}

			user, err := userService.GetUserByUsername(username)

			if err != nil || subtle.ConstantTimeCompare([]byte(password), []byte(user.Password)) != 1 {
				basicAuthFailed(w)
				return
			}

			ctx := context.WithValue(r.Context(), userKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func basicAuthFailed(w http.ResponseWriter) {
	w.Header().Add("WWW-Authenticate", `Basic realm="Tournament"`)
	w.WriteHeader(http.StatusUnauthorized)
}
