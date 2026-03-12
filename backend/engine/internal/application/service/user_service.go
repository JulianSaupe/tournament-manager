package service

import (
	"context"
	"engine/internal/domain"
	"engine/internal/ports/input"
	"engine/internal/ports/output"
)

// UserService implements the UserService interface
type UserService struct {
	userRepository output.UserRepositoryInterface
}

// NewUserService creates a new user service
func NewUserService(userRepository output.UserRepositoryInterface) input.UserServiceInterface {
	return &UserService{
		userRepository: userRepository,
	}
}

// GetUserByUsername retrieves a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	return s.userRepository.FindByUsername(ctx, username)
}
