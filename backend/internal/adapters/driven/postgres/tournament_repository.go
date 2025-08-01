package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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

	db.RegisterModel((*domain.Tournament)(nil))
	db.NewCreateTable().Model((*domain.Tournament)(nil)).IfNotExists().Exec(context.Background())

	return &TournamentRepository{
		db: db,
	}, nil
}

// FindByID retrieves a tournament by its Id
func (r *TournamentRepository) FindByID(ctx context.Context, id string) (*domain.Tournament, error) {
	// Add timeout for database operations
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	tournament := new(domain.Tournament)

	err := r.db.NewSelect().
		Model(tournament).
		Where("id = ?", id).
		Column("tournament.*").
		ColumnExpr(`(
            SELECT count(*)
            FROM player p
            WHERE p.tournament_id = tournament.id
        ) AS player_count`).
		Relation("players").
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("tournament not found")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out: %w", err)
		}
		return nil, fmt.Errorf("error finding tournament: %w", err)
	}

	return tournament, nil
}

// FindAll retrieves all tournaments
func (r *TournamentRepository) FindAll(ctx context.Context) ([]*domain.IndexTournament, error) {
	tournaments := make([]*domain.IndexTournament, 0)

	err := r.db.NewSelect().
		Model(&tournaments).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error querying tournaments: %w", err)
	}

	return tournaments, nil
}

// Save persists a tournament
func (r *TournamentRepository) Save(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error) {
	_, err := r.db.NewInsert().Model(tournament).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("error saving tournament: %w", err)
	}

	return tournament, nil
}

// Delete removes a tournament
func (r *TournamentRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.NewDelete().
		Model((*domain.Tournament)(nil)).
		Where("id = ?", id).
		Exec(ctx)

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

func (r *TournamentRepository) Update(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error) {
	err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		result, err := tx.NewUpdate().
			Model(tournament).
			Set("status = ?", tournament.Status).
			WherePK().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("error updating tournament: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return domain.NewNotFoundError("tournament not found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return tournament, nil
}
