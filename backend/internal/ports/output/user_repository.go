package output

import (
	"Tournament/internal/domain"
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByUsername retrieves a user by their username
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
}
