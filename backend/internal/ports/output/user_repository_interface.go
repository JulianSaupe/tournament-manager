package output

import (
	"Tournament/internal/domain"
	"context"
)

// UserRepositoryInterface defines the interface for user data access
type UserRepositoryInterface interface {
	// FindByUsername retrieves a user by their username
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}
