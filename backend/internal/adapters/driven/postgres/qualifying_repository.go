package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
)

type QualifyingRepository struct {
	db *sql.DB
}

func NewQualifyingRepository(db *sql.DB) (output.QualifyingRepositoryInterface, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}
	return &QualifyingRepository{
		db: db,
	}, nil
}

func (r *QualifyingRepository) FindByTournamentId(ctx context.Context, id string) (*domain.Qualifying, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	query := `
			SELECT q.player_id, p.name, RANK() OVER (ORDER BY q.time) as position,  q.created_at, q.time
			FROM qualifying q JOIN players p ON q.player_id = p.id
			WHERE q.tournament_id = $1
		`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("error querying tournaments: %w", err)
	}
	defer r.closeRows(rows)

	players := make([]*domain.QualifyingPlayer, 0)
	for rows.Next() {
		player := new(domain.QualifyingPlayer)
		err := rows.Scan(
			&player.PlayerId,
			&player.PlayerName,
			&player.Position,
			&player.SignupDate,
			&player.BestTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tournament: %w", err)
		}
		players = append(players, player)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tournaments: %w", err)
	}

	qualifying := domain.Qualifying{
		TournamentId: id,
		Players:      players,
	}

	return &qualifying, nil
}

func (r *QualifyingRepository) DeleteByTournamentId(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	query := `DELETE FROM qualifying WHERE tournament_id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting tournament: %w", err)
	}

	return r.checkRowsAffected(result, "qualifying not found")
}

func (r *QualifyingRepository) checkRowsAffected(result sql.Result, notFoundMsg string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.NewNotFoundError(notFoundMsg)
	}
	return nil
}

func (r *QualifyingRepository) closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
	}
}
