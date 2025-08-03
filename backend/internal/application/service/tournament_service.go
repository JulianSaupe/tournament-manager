package service

import (
	"Tournament/internal/adapters/driving/requests"
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
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
func (s *TournamentService) CreateTournament(ctx context.Context, req *requests.CreateTournamentRequest) *domain.Tournament {
	rounds := make([]domain.Round, 0)

	for _, round := range req.Rounds {
		var newRound = domain.Round{
			Name:                   round.Name,
			MatchCount:             round.MatchCount,
			PlayerAdvancementCount: round.PlayerAdvancementCount,
			PlayerCount:            round.GroupCount * round.GroupSize,
			GroupSize:              round.GroupSize,
			ConcurrentGroupCount:   round.ConcurrentGroupCount,
			Groups:                 make([]domain.Group, 0),
		}

		rounds = append(rounds, newRound)
	}

	var newTournament = domain.Tournament{
		Name:                   req.Name,
		Description:            req.Description,
		StartDate:              req.StartDate,
		EndDate:                req.EndDate,
		Status:                 domain.StatusDraft,
		Rounds:                 rounds,
		AllowUnderfilledGroups: req.AllowUnderfilledGroups,
	}

	savedTournament, err := s.tournamentRepository.InsertNewTournament(ctx, &newTournament)

	if err != nil {
		panic(err)
	}

	return savedTournament
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
