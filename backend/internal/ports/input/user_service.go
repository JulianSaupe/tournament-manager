package input

import (
	"Tournament/internal/domain"
)

// UserService defines the interface for user business operations
type UserService interface {
	// GetUserByUsername retrieves a user by username
	GetUserByUsername(username string) (*domain.User, error)
}
