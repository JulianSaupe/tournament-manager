package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

// TournamentRepository is a PostgreSQL implementation of the TournamentRepository interface
type TournamentRepository struct {
	db *sql.DB
}

// NewPostgresTournamentRepository creates a new PostgreSQL tournament repository
func NewPostgresTournamentRepository(connectionString string) (output.TournamentRepository, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &TournamentRepository{
		db: db,
	}, nil
}

// FindByID retrieves a tournament by its ID
func (r *TournamentRepository) FindByID(id string) (*domain.Tournament, error) {
	query := `
	SELECT id, name, description, start_date, end_date, status
	FROM tournaments
	WHERE id = $1;`

	var tournament domain.Tournament
	err := r.db.QueryRow(query, id).Scan(
		&tournament.ID,
		&tournament.Name,
		&tournament.Description,
		&tournament.StartDate,
		&tournament.EndDate,
		&tournament.Status,
	)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("tournament not found")
	}

	if err != nil {
		return nil, fmt.Errorf("error finding tournament: %w", err)
	}

	return &tournament, nil
}

// FindAll retrieves all tournaments
func (r *TournamentRepository) FindAll() ([]*domain.Tournament, error) {
	query := `
	SELECT id, name, description, start_date, end_date, status
	FROM tournaments;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying tournaments: %w", err)
	}
	defer rows.Close()

	tournaments := make([]*domain.Tournament, 0)
	for rows.Next() {
		var tournament domain.Tournament
		err := rows.Scan(
			&tournament.ID,
			&tournament.Name,
			&tournament.Description,
			&tournament.StartDate,
			&tournament.EndDate,
			&tournament.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning tournament row: %w", err)
		}
		tournaments = append(tournaments, &tournament)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tournament rows: %w", err)
	}

	return tournaments, nil
}

// Save persists a tournament
func (r *TournamentRepository) Save(tournament *domain.Tournament) (*domain.Tournament, error) {
	query := `
	INSERT INTO tournaments (id, name, description, start_date, end_date, status)
	VALUES ($1, $2, $3, $4, $5, $6)
	ON CONFLICT (id) DO UPDATE
	SET name = $2, description = $3, start_date = $4, end_date = $5, status = $6
	RETURNING id, name, description, start_date, end_date, status;`

	var savedTournament domain.Tournament
	err := r.db.QueryRow(
		query,
		tournament.ID,
		tournament.Name,
		tournament.Description,
		tournament.StartDate,
		tournament.EndDate,
		tournament.Status,
	).Scan(
		&savedTournament.ID,
		&savedTournament.Name,
		&savedTournament.Description,
		&savedTournament.StartDate,
		&savedTournament.EndDate,
		&savedTournament.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("error saving tournament: %w", err)
	}

	return &savedTournament, nil
}

// Delete removes a tournament
func (r *TournamentRepository) Delete(id string) error {
	query := `DELETE FROM tournaments WHERE id = $1;`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting tournament: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("tournament not found")
	}

	return nil
}

// Close closes the database connection
func (r *TournamentRepository) Close() error {
	return r.db.Close()
}
