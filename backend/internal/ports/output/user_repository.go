package output

import (
	"Tournament/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByUsername retrieves a user by their username
	FindByUsername(username string) (*domain.User, error)
}
