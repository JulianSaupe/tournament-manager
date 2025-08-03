package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// UserRepository is a PostgreSQL implementation of the UserRepository interface
type UserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *sql.DB) (output.UserRepositoryInterface, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	return &UserRepository{
		db: db,
	}, nil
}

// FindByUsername retrieves a user by their username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := `
		SELECT id, username, password, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	row := r.db.QueryRowContext(ctx, query, username)

	user := new(domain.User)
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&createdAt,
		&updatedAt,
	)

	user.CreatedAt = createdAt
	user.UpdatedAt = updatedAt

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}
