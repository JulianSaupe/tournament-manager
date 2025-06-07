package output

import (
	"Tournament/internal/domain"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// FindByUsername retrieves a user by their username
	FindByUsername(username string) (*domain.User, error)
	
	// Save persists a user
	Save(user *domain.User) (*domain.User, error)
	
	// FindAll retrieves all users
	FindAll() ([]*domain.User, error)
	
	// Delete removes a user
	Delete(id int) error
}