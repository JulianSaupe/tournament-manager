package input

import (
	"Tournament/internal/domain"
	"context"
)

// UserService defines the interface for user business operations
type UserService interface {
	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}
