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

type TournamentRepository struct {
	db *sql.DB
}

// NewTournamentRepository creates a new PostgreSQL tournament repository
func NewTournamentRepository(db *sql.DB) (output.TournamentRepository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	return &TournamentRepository{
		db: db,
	}, nil
}

// FindByID retrieves a tournament by its Id
func (r *TournamentRepository) FindByID(ctx context.Context, id string) (*domain.Tournament, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Query tournament
	query := `
		SELECT id, name, description, start_date, end_date, status, player_count
		FROM tournaments
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	tournament := new(domain.Tournament)
	err := row.Scan(
		&tournament.Id,
		&tournament.Name,
		&tournament.Description,
		&tournament.StartDate,
		&tournament.EndDate,
		&tournament.Status,
		&tournament.PlayerCount,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("tournament not found")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out: %w", err)
		}
		return nil, fmt.Errorf("error finding tournament: %w", err)
	}

	// Query players for this tournament
	playersQuery := `
		SELECT id, name, tournament_id
		FROM players
		WHERE tournament_id = $1
	`
	rows, err := r.db.QueryContext(ctx, playersQuery, id)
	if err != nil {
		return nil, fmt.Errorf("error querying players: %w", err)
	}
	defer rows.Close()

	players := []domain.Player{}
	for rows.Next() {
		player := domain.Player{}
		err := rows.Scan(&player.Id, &player.Name, &player.TournamentId)
		if err != nil {
			return nil, fmt.Errorf("error scanning player: %w", err)
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating players: %w", err)
	}

	tournament.Players = players

	return tournament, nil
}

// FindAll retrieves all tournaments
func (r *TournamentRepository) FindAll(ctx context.Context) ([]*domain.IndexTournament, error) {
	query := `
		SELECT id, name, description, start_date, end_date, status
		FROM tournaments
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying tournaments: %w", err)
	}
	defer rows.Close()

	tournaments := make([]*domain.IndexTournament, 0)
	for rows.Next() {
		tournament := new(domain.IndexTournament)
		err := rows.Scan(
			&tournament.Id,
			&tournament.Name,
			&tournament.Description,
			&tournament.StartDate,
			&tournament.EndDate,
			&tournament.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tournament: %w", err)
		}
		tournaments = append(tournaments, tournament)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tournaments: %w", err)
	}

	return tournaments, nil
}

// Save persists a tournament
func (r *TournamentRepository) Save(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error) {
	query := `
		INSERT INTO tournaments (id, name, description, start_date, end_date, status, player_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		tournament.Id,
		tournament.Name,
		tournament.Description,
		tournament.StartDate,
		tournament.EndDate,
		tournament.Status,
		tournament.PlayerCount,
	)

	if err != nil {
		return nil, fmt.Errorf("error saving tournament: %w", err)
	}

	return tournament, nil
}

// Delete removes a tournament
func (r *TournamentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tournaments WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)

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
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	query := `
		UPDATE tournaments
		SET status = $1
		WHERE id = $2
	`
	result, err := tx.ExecContext(ctx, query, tournament.Status, tournament.Id)
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

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return tournament, nil
}
