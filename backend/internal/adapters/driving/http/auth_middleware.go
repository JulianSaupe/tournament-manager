package http

import (
	"Tournament/internal/ports/input"
	"context"
	"net/http"
	"strings"
)

// userKey is the context key for the authenticated user
type userKey struct{}

// AuthMiddleware creates a middleware that authenticates requests
func AuthMiddleware(userService input.UserService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get the Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				SendErrorResponse(w, r, http.StatusUnauthorized, "Authorization header required")
				return
			}

			// Check if the header has the correct format
			if !strings.HasPrefix(authHeader, "Basic ") {
				SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid authorization format")
				return
			}

			// Extract the token
			token := strings.TrimPrefix(authHeader, "Basic ")
			if token == "" {
				SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Split the token to get username and timestamp
			parts := strings.Split(token, ":")
			if len(parts) != 2 {
				SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid token format")
				return
			}

			username := parts[0]

			// Get the user from the database
			user, err := userService.GetUserByUsername(username)
			if err != nil {
				SendErrorResponse(w, r, http.StatusUnauthorized, "Invalid token")
				return
			}

			// Add the user to the request context
			ctx := context.WithValue(r.Context(), userKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
