package application

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"time"
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

// RegisterUser registers a new user
func (s *UserServiceImpl) RegisterUser(username, password string) (*domain.User, error) {
	// Check if username already exists
	existingUser, err := s.userRepository.FindByUsername(username)
	if err == nil && existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	user := &domain.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return s.userRepository.Save(user)
}

// AuthenticateUser authenticates a user with username and password
func (s *UserServiceImpl) AuthenticateUser(username, password string) (string, error) {
	user, err := s.userRepository.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Compare the provided password with the stored hash
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// For simplicity, we're just returning a basic token
	// In a real application, you would use JWT or another token mechanism
	token := fmt.Sprintf("%s:%d", username, time.Now().Unix())
	return token, nil
}

// GetUserByUsername retrieves a user by username
func (s *UserServiceImpl) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepository.FindByUsername(username)
}

// ListUsers retrieves all users
func (s *UserServiceImpl) ListUsers() ([]*domain.User, error) {
	return s.userRepository.FindAll()
}