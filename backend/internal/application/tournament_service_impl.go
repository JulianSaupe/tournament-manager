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
func (s *TournamentServiceImpl) CreateTournament(name, description, startDate, endDate string) *domain.Tournament {
	tournament := &domain.Tournament{
		ID:          uuid.New().String(),
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      domain.StatusDraft,
	}

	tournament, err := s.tournamentRepository.Save(tournament)

	if err != nil {
		panic(err)
	}

	return tournament
}

// GetTournament retrieves a tournament by ID
func (s *TournamentServiceImpl) GetTournament(id string) *domain.Tournament {
	tournament, err := s.tournamentRepository.FindByID(id)

	if err != nil {
		panic(err)
	}

	return tournament
}

// ListTournaments retrieves all tournaments
func (s *TournamentServiceImpl) ListTournaments() []*domain.Tournament {
	tournaments, err := s.tournamentRepository.FindAll()

	if err != nil {
		panic(err)
	}

	return tournaments
}

// UpdateTournamentStatus updates the status of a tournament
func (s *TournamentServiceImpl) UpdateTournamentStatus(id string, status domain.TournamentStatus) *domain.Tournament {
	tournament, err := s.tournamentRepository.FindByID(id)
	if err != nil {
		panic(err)
	}

	tournament.Status = status
	tournament, err = s.tournamentRepository.Update(tournament)

	if err != nil {
		panic(err)
	}

	return tournament
}

func (s *TournamentServiceImpl) DeleteTournament(id string) {
	err := s.tournamentRepository.Delete(id)

	if err != nil {
		panic(err)
	}
}
