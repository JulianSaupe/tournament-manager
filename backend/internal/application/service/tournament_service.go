package service

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
	"github.com/google/uuid"
)

// TournamentService implements the TournamentService interface
type TournamentService struct {
	tournamentRepository output.TournamentRepository
}

// NewTournamentService creates a new tournament service
func NewTournamentService(tournamentRepository output.TournamentRepository) input.TournamentService {
	return &TournamentService{
		tournamentRepository: tournamentRepository,
	}
}

// CreateTournament creates a new tournament
func (s *TournamentService) CreateTournament(ctx context.Context, name, description, startDate, endDate string) *domain.Tournament {
	tournament := &domain.Tournament{
		Id:          uuid.New().String(),
		Name:        name,
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      domain.StatusDraft,
	}

	tournament, err := s.tournamentRepository.Save(ctx, tournament)

	if err != nil {
		panic(err)
	}

	return tournament
}

// GetTournament retrieves a tournament by Id
func (s *TournamentService) GetTournament(ctx context.Context, id string) *domain.Tournament {
	tournament, err := s.tournamentRepository.FindByID(ctx, id)

	if err != nil {
		panic(err)
	}

	return tournament
}

// ListTournaments retrieves all tournaments
func (s *TournamentService) ListTournaments(ctx context.Context) []*domain.IndexTournament {
	tournaments, err := s.tournamentRepository.FindAll(ctx)

	if err != nil {
		panic(err)
	}

	return tournaments
}

// UpdateTournamentStatus updates the status of a tournament
func (s *TournamentService) UpdateTournamentStatus(ctx context.Context, id string, status domain.TournamentStatus) *domain.Tournament {
	tournament, err := s.tournamentRepository.FindByID(ctx, id)
	if err != nil {
		panic(err)
	}

	tournament.Status = status
	tournament, err = s.tournamentRepository.Update(ctx, tournament)

	if err != nil {
		panic(err)
	}

	return tournament
}

func (s *TournamentService) DeleteTournament(ctx context.Context, id string) {
	err := s.tournamentRepository.Delete(ctx, id)

	if err != nil {
		panic(err)
	}
}
