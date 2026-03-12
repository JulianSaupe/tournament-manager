package middleware

import (
	"engine/internal/application/service"
	"engine/internal/domain"
	"net/http"
)

// AuthorizationMiddleware checks if the user has permission to access a resource
func AuthorizationMiddleware(authorizationService *service.AuthorizationService, resource string, action string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := GetUserIDFromContext(r.Context())
			if !ok {
				panic(domain.NewUnauthorizedError("User not authenticated"))
			}

			ctx := r.Context()
			allowed, message, err := authorizationService.CheckPermission(ctx, userID, resource, action)
			if err != nil {
				panic(domain.NewForbiddenError("Failed to check permission: " + err.Error()))
			}

			if !allowed {
				panic(domain.NewForbiddenError("Permission denied: " + message))
			}

			next.ServeHTTP(w, r)
		})
	}
}
