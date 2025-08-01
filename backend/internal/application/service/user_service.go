package service

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
)

// UserService implements the UserService interface
type UserService struct {
	userRepository output.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepository output.UserRepository) input.UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}
