package application

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"github.com/google/uuid"
)

// TournamentServiceImpl implements the TournamentService interface
type TournamentServiceImpl struct {
	tournamentRepository output.TournamentRepository
}

// NewTournamentService creates a new tournament service
func NewTournamentService(tournamentRepository output.TournamentRepository) input.TournamentService {
	return &TournamentServiceImpl{
		tournamentRepository: tournamentRepository,
	}
}

// CreateTournament creates a new tournament
func (s *TournamentServiceImpl) CreateTournament(name, description, startDate, endDate string) (*domain.Tournament, error) {
	tournament := &domain.Tournament{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      domain.StatusDraft,
	}

	return s.tournamentRepository.Save(tournament)
}

// GetTournament retrieves a tournament by ID
func (s *TournamentServiceImpl) GetTournament(id string) (*domain.Tournament, error) {
	return s.tournamentRepository.FindByID(id)
}

// ListTournaments retrieves all tournaments
func (s *TournamentServiceImpl) ListTournaments() ([]*domain.Tournament, error) {
	return s.tournamentRepository.FindAll()
}

// UpdateTournamentStatus updates the status of a tournament
func (s *TournamentServiceImpl) UpdateTournamentStatus(id string, status domain.TournamentStatus) (*domain.Tournament, error) {
	tournament, err := s.tournamentRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	tournament.Status = status
	return s.tournamentRepository.Update(tournament)
}
