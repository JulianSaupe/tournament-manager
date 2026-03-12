package middleware

import (
	"context"
	"engine/internal/ports/input"
	"net/http"
)

type userIDKey struct{}

// AuthenticationMiddleware validates the session cookie and stores user_id in context
func AuthenticationMiddleware(authService input.AuthenticationServiceInterface) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get session_id from cookie
			cookie, err := r.Cookie("session_id")
			if err != nil || cookie.Value == "" {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized: missing session"))
				return
			}

			// Validate session via gRPC
			ctx := r.Context()
			userID, err := authService.ValidateSession(ctx, cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized: invalid session"))
				return
			}

			// Store user_id in context
			ctx = context.WithValue(ctx, userIDKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserIDFromContext retrieves the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(userIDKey{}).(string)
	return userID, ok
}
