package application

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
)

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepository output.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepository output.UserRepository) input.UserService {
	return &UserServiceImpl{
		userRepository: userRepository,
	}
}

// GetUserByUsername retrieves a user by username
func (s *UserServiceImpl) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}
