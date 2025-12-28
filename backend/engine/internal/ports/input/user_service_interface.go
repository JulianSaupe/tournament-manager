package input

import (
	"Tournament/internal/domain"
	"context"
)

// UserServiceInterface defines the interface for user business operations
type UserServiceInterface interface {
	// GetUserByUsername retrieves a user by username
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
}
