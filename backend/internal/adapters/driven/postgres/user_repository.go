package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
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
func NewPostgresUserRepository(connectionString string) (output.UserRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &UserRepository{
		db: db,
	}, nil
}

// FindByUsername retrieves a user by their username
func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	query := `
	SELECT id, username, password, created_at, updated_at
	FROM users
	WHERE username = $1;`

	var user domain.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}

	return &user, nil
}

// Save persists a user
func (r *UserRepository) Save(user *domain.User) (*domain.User, error) {
	// For new users (ID = 0), insert a new record
	if user.ID == 0 {
		query := `
		INSERT INTO users (username, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, password, created_at, updated_at;`

		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now

		var savedUser domain.User
		err := r.db.QueryRow(
			query,
			user.Username,
			user.Password,
			user.CreatedAt,
			user.UpdatedAt,
		).Scan(
			&savedUser.ID,
			&savedUser.Username,
			&savedUser.Password,
			&savedUser.CreatedAt,
			&savedUser.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error saving user: %w", err)
		}

		return &savedUser, nil
	}

	// For existing users, update the record
	query := `
	UPDATE users
	SET username = $1, password = $2, updated_at = $3
	WHERE id = $4
	RETURNING id, username, password, created_at, updated_at;`

	user.UpdatedAt = time.Now()

	var updatedUser domain.User
	err := r.db.QueryRow(
		query,
		user.Username,
		user.Password,
		user.UpdatedAt,
		user.ID,
	).Scan(
		&updatedUser.ID,
		&updatedUser.Username,
		&updatedUser.Password,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return &updatedUser, nil
}

// FindAll retrieves all users
func (r *UserRepository) FindAll() ([]*domain.User, error) {
	query := `
	SELECT id, username, password, created_at, updated_at
	FROM users;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Password,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user row: %w", err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// Delete removes a user
func (r *UserRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1;`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

// Close closes the database connection
func (r *UserRepository) Close() error {
	return r.db.Close()
}