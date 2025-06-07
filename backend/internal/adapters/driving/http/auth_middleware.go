package http

import (
	"Tournament/internal/ports/input"
	"context"
	"golang.org/x/crypto/bcrypt"
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

			if err != nil || checkPasswordHash(password, user.Password) != true {
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

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}
