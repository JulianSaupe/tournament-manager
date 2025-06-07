package memory

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"errors"
	"github.com/google/uuid"
	"sync"
)

// InMemoryTournamentRepository is an in-memory implementation of the TournamentRepository interface
type InMemoryTournamentRepository struct {
	tournaments map[string]*domain.Tournament
	mutex       sync.RWMutex
}

// NewInMemoryTournamentRepository creates a new in-memory tournament repository
func NewInMemoryTournamentRepository() output.TournamentRepository {
	return &InMemoryTournamentRepository{
		tournaments: make(map[string]*domain.Tournament),
	}
}

// FindByID retrieves a tournament by its ID
func (r *InMemoryTournamentRepository) FindByID(id string) (*domain.Tournament, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tournament, exists := r.tournaments[id]
	if !exists {
		return nil, errors.New("tournament not found")
	}

	return tournament, nil
}

// FindAll retrieves all tournaments
func (r *InMemoryTournamentRepository) FindAll() ([]*domain.Tournament, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tournaments := make([]*domain.Tournament, 0, len(r.tournaments))
	for _, tournament := range r.tournaments {
		tournaments = append(tournaments, tournament)
	}

	return tournaments, nil
}

// Save persists a tournament
func (r *InMemoryTournamentRepository) Save(tournament *domain.Tournament) (*domain.Tournament, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// If the tournament doesn't have an ID, assign one
	if tournament.ID == "" {
		tournament.ID = generateID()
	}

	// Create a copy to avoid modifying the original
	tournamentCopy := &domain.Tournament{
		ID:          tournament.ID,
		Name:        tournament.Name,
		Description: tournament.Description,
		StartDate:   tournament.StartDate,
		EndDate:     tournament.EndDate,
		Status:      tournament.Status,
	}

	r.tournaments[tournament.ID] = tournamentCopy
	return tournamentCopy, nil
}

// Delete removes a tournament
func (r *InMemoryTournamentRepository) Delete(id string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.tournaments[id]; !exists {
		return errors.New("tournament not found")
	}

	delete(r.tournaments, id)
	return nil
}

// generateID generates a UUID string
func generateID() string {
	return uuid.New().String()
}
