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

// TournamentRepository is a PostgreSQL implementation of the TournamentRepository interface
type TournamentRepository struct {
	db *bun.DB
}

// NewTournamentRepository creates a new PostgreSQL tournament repository
func NewTournamentRepository(db *bun.DB) (output.TournamentRepository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	// Register models with Bun
	db.RegisterModel((*domain.Tournament)(nil))

	return &TournamentRepository{
		db: db,
	}, nil
}

// FindByID retrieves a tournament by its Id
func (r *TournamentRepository) FindByID(id string) (*domain.Tournament, error) {
	tournament := new(domain.Tournament)

	err := r.db.NewSelect().
		Model(tournament).
		Where("id = ?", id).
		Scan(context.Background())

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("tournament not found")
		}
		return nil, fmt.Errorf("error finding tournament: %w", err)
	}

	return tournament, nil
}

// FindAll retrieves all tournaments
func (r *TournamentRepository) FindAll() ([]*domain.Tournament, error) {
	tournaments := make([]*domain.Tournament, 0)

	err := r.db.NewSelect().
		Model(&tournaments).
		Scan(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error querying tournaments: %w", err)
	}

	return tournaments, nil
}

// Save persists a tournament
func (r *TournamentRepository) Save(tournament *domain.Tournament) (*domain.Tournament, error) {
	_, err := r.db.NewInsert().Model(tournament).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error saving tournament: %w", err)
	}

	return tournament, nil
}

// Delete removes a tournament
func (r *TournamentRepository) Delete(id string) error {
	result, err := r.db.NewDelete().
		Model((*domain.Tournament)(nil)).
		Where("id = ?", id).
		Exec(context.Background())

	if err != nil {
		return fmt.Errorf("error deleting tournament: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("tournament not found")
	}

	return nil
}

func (r *TournamentRepository) Update(tournament *domain.Tournament) (*domain.Tournament, error) {
	result, err := r.db.NewUpdate().
		Model(tournament).
		Set("status = ?", tournament.Status).
		WherePK().
		Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error updating tournament: %w", err)
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, domain.NewNotFoundError("tournament not found")
	}

	return tournament, nil
}
