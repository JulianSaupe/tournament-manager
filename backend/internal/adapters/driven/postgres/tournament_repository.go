package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	queryTimeout = 30 * time.Second
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
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	tournament, err := r.findTournamentByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return r.populateTournamentDetails(ctx, tournament)
}

// FindAll retrieves all tournaments
func (r *TournamentRepository) FindAll(ctx context.Context) ([]*domain.IndexTournament, error) {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	query := `
		SELECT id, name, description, start_date, end_date, status
		FROM tournaments
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying tournaments: %w", err)
	}
	defer r.closeRows(rows)

	return r.scanTournamentList(rows)
}

func (r *TournamentRepository) InsertNewTournament(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error) {
	return r.executeInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) (*domain.Tournament, error) {
		tournamentID, err := r.insertTournament(ctx, tx, tournament)
		if err != nil {
			return nil, err
		}
		tournament.Id = tournamentID

		if len(tournament.Rounds) > 0 {
			err = r.insertRounds(ctx, tx, tournament.Rounds, tournamentID)
			if err != nil {
				return nil, err
			}
		}

		return tournament, nil
	})
}

// Delete removes a tournament
func (r *TournamentRepository) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	query := `DELETE FROM tournaments WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("error deleting tournament: %w", err)
	}

	return r.checkRowsAffected(result, "tournament not found")
}

func (r *TournamentRepository) Update(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error) {
	return r.executeInTransaction(ctx, func(ctx context.Context, tx *sql.Tx) (*domain.Tournament, error) {
		query := `
			UPDATE tournaments
			SET status = $1
			WHERE id = $2
		`
		result, err := tx.ExecContext(ctx, query, tournament.Status, tournament.Id)
		if err != nil {
			return nil, fmt.Errorf("error updating tournament: %w", err)
		}

		err = r.checkRowsAffected(result, "tournament not found")
		if err != nil {
			return nil, err
		}

		return tournament, nil
	})
}

// Helper methods

func (r *TournamentRepository) executeInTransaction(ctx context.Context, fn func(context.Context, *sql.Tx) (*domain.Tournament, error)) (*domain.Tournament, error) {
	tx, err := r.db.BeginTx(ctx, &sql.TxOptions{})
	if err != nil {
		return nil, fmt.Errorf("error starting transaction: %w", err)
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				err = errors.Join(err, fmt.Errorf("rollback failed: %w", rollbackErr))
			}
		}
	}()

	result, err := fn(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transaction: %w", err)
	}

	return result, nil
}

func (r *TournamentRepository) populateTournamentDetails(ctx context.Context, tournament *domain.Tournament) (*domain.Tournament, error) {
	players, err := r.findPlayersByTournamentID(ctx, tournament.Id)
	if err != nil {
		return nil, err
	}
	tournament.Players = players

	rounds, err := r.findRoundsByTournamentID(ctx, tournament.Id)
	if err != nil {
		return nil, err
	}
	tournament.Rounds = rounds

	return tournament, nil
}

func (r *TournamentRepository) scanTournamentList(rows *sql.Rows) ([]*domain.IndexTournament, error) {
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

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tournaments: %w", err)
	}

	return tournaments, nil
}

func (r *TournamentRepository) checkRowsAffected(result sql.Result, notFoundMsg string) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}
	if rowsAffected == 0 {
		return domain.NewNotFoundError(notFoundMsg)
	}
	return nil
}

func (r *TournamentRepository) closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Printf("failed to close rows: %v", err)
	}
}

func (r *TournamentRepository) findTournamentByID(ctx context.Context, id string) (*domain.Tournament, error) {
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
	return tournament, nil
}

func (r *TournamentRepository) findPlayersByTournamentID(ctx context.Context, tournamentID string) ([]domain.Player, error) {
	query := `
		SELECT id, name, tournament_id
		FROM players
		WHERE tournament_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, tournamentID)
	if err != nil {
		return nil, fmt.Errorf("error querying players: %w", err)
	}
	defer r.closeRows(rows)

	var players []domain.Player
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
	return players, nil
}

func (r *TournamentRepository) findRoundsByTournamentID(ctx context.Context, tournamentID string) ([]domain.Round, error) {
	query := `
		SELECT id, name, match_count, player_count, player_advancement_count, group_size, concurrent_group_count
		FROM rounds
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error querying rounds: %w", err)
	}
	defer r.closeRows(rows)

	var rounds []domain.Round
	for rows.Next() {
		round := domain.Round{}
		err := rows.Scan(&round.Id, &round.Name, &round.MatchCount, &round.PlayerCount, &round.PlayerAdvancementCount, &round.GroupSize, &round.ConcurrentGroupCount)
		if err != nil {
			return nil, fmt.Errorf("error scanning round: %w", err)
		}
		round.TournamentId = tournamentID
		round.Groups = make([]domain.Group, 0)
		rounds = append(rounds, round)
	}
	return rounds, nil
}

func (r *TournamentRepository) insertTournament(ctx context.Context, tx *sql.Tx, tournament *domain.Tournament) (string, error) {
	var tournamentID string
	query := `
        INSERT INTO tournaments (id, name, description, start_date, end_date, status, player_count, allow_underfilled_groups)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `
	err := tx.QueryRowContext(
		ctx,
		query,
		tournament.Name,
		tournament.Description,
		tournament.StartDate,
		tournament.EndDate,
		tournament.Status,
		tournament.PlayerCount,
		tournament.AllowUnderfilledGroups,
	).Scan(&tournamentID)
	if err != nil {
		return "", fmt.Errorf("error saving tournament: %w", err)
	}
	return tournamentID, nil
}

func (r *TournamentRepository) insertRounds(ctx context.Context, tx *sql.Tx, rounds []domain.Round, tournamentID string) error {
	placeholders := make([]string, len(rounds))
	args := make([]interface{}, 0, len(rounds)*7)
	for i, round := range rounds {
		start := i*7 + 1
		placeholders[i] = fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d)",
			start, start+1, start+2, start+3, start+4, start+5, start+6)
		args = append(args,
			round.Name,
			tournamentID,
			round.MatchCount,
			round.PlayerCount,
			round.PlayerAdvancementCount,
			round.GroupSize,
			round.ConcurrentGroupCount,
		)
	}
	roundQuery := fmt.Sprintf(`
        INSERT INTO rounds (name, tournament_id, match_count, player_count, player_advancement_count, group_size, concurrent_group_count)
        VALUES %s`, strings.Join(placeholders, ", "))
	_, err := tx.ExecContext(ctx, roundQuery, args...)
	if err != nil {
		return fmt.Errorf("error saving rounds: %w", err)
	}
	return nil
}
