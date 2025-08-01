package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

// UserRepository is a PostgreSQL implementation of the UserRepository interface
type UserRepository struct {
	db *bun.DB
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db *bun.DB) (output.UserRepository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	db.RegisterModel((*domain.User)(nil))
	db.NewCreateTable().Model((*domain.User)(nil)).IfNotExists().Exec(context.Background())

	return &UserRepository{
		db: db,
	}, nil
}

// FindByUsername retrieves a user by their username
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	user := new(domain.User)

	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return user, nil
}
