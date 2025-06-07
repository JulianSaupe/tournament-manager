package input

import (
	"Tournament/internal/domain"
)

// UserService defines the interface for user business operations
type UserService interface {
	// RegisterUser registers a new user
	RegisterUser(username, password string) (*domain.User, error)
	
	// AuthenticateUser authenticates a user with username and password
	AuthenticateUser(username, password string) (string, error)
	
	// GetUserByUsername retrieves a user by username
	GetUserByUsername(username string) (*domain.User, error)
	
	// ListUsers retrieves all users
	ListUsers() ([]*domain.User, error)
}