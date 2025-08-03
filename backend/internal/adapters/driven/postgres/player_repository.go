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

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) (output.PlayerRepositoryInterface, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	return &PlayerRepository{
		db: db,
	}, nil
}

func (r *PlayerRepository) InsertNewPlayer(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	query := `
		INSERT INTO players (name, tournament_id)
		VALUES ($1, $2)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		player.Name,
		player.TournamentId,
	)

	if err != nil {
		return nil, fmt.Errorf("error saving player: %w", err)
	}

	return player, nil
}

func (r *PlayerRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM players WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)

	if err != nil {
		return fmt.Errorf("error deleting player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("player not found")
	}

	return nil
}

func (r *PlayerRepository) FindAll(ctx context.Context, tournamentId string) ([]*domain.Player, error) {
	query := `
		SELECT id, name, tournament_id
		FROM players
		WHERE tournament_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, tournamentId)
	if err != nil {
		return nil, fmt.Errorf("error querying players: %w", err)
	}
	defer rows.Close()

	players := make([]*domain.Player, 0)
	for rows.Next() {
		player := new(domain.Player)
		err := rows.Scan(
			&player.Id,
			&player.Name,
			&player.TournamentId,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning player: %w", err)
		}
		players = append(players, player)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating players: %w", err)
	}

	return players, nil
}

func (r *PlayerRepository) FindByID(ctx context.Context, id string) (*domain.Player, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	query := `
		SELECT id, name, tournament_id
		FROM players
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	player := new(domain.Player)
	err := row.Scan(
		&player.Id,
		&player.Name,
		&player.TournamentId,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("player not found")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out: %w", err)
		}
		return nil, fmt.Errorf("error finding player: %w", err)
	}

	return player, nil
}

func (r *PlayerRepository) UpdateName(ctx context.Context, player *domain.Player) (*domain.Player, error) {
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
		UPDATE players
		SET name = $1
		WHERE id = $2
	`
	result, err := tx.ExecContext(ctx, query, player.Name, player.Id)
	if err != nil {
		return nil, fmt.Errorf("error updating player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, domain.NewNotFoundError("player not found")
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return player, nil
}
