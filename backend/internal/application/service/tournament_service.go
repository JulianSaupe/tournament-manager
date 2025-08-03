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
	newTournament := s.buildTournamentFromRequest(req)
	savedTournament, err := s.tournamentRepository.InsertNewTournament(ctx, &newTournament)
	s.handleRepositoryError(err)
	return savedTournament
}

// GetTournament retrieves a tournament by Id
func (s *TournamentService) GetTournament(ctx context.Context, id string) *domain.Tournament {
	tournament, err := s.tournamentRepository.FindByID(ctx, id)
	s.handleRepositoryError(err)
	return tournament
}

// ListTournaments retrieves all tournaments
func (s *TournamentService) ListTournaments(ctx context.Context) []*domain.IndexTournament {
	tournaments, err := s.tournamentRepository.FindAll(ctx)
	s.handleRepositoryError(err)
	return tournaments
}

// UpdateTournamentStatus updates the status of a tournament
func (s *TournamentService) UpdateTournamentStatus(ctx context.Context, id string, status domain.TournamentStatus) *domain.Tournament {
	tournament, err := s.tournamentRepository.FindByID(ctx, id)
	s.handleRepositoryError(err)

	tournament.Status = status
	tournament, err = s.tournamentRepository.Update(ctx, tournament)
	s.handleRepositoryError(err)
	return tournament
}

func (s *TournamentService) DeleteTournament(ctx context.Context, id string) {
	err := s.tournamentRepository.Delete(ctx, id)
	s.handleRepositoryError(err)
}

// buildTournamentFromRequest constructs a domain Tournament from a request
func (s *TournamentService) buildTournamentFromRequest(req *requests.CreateTournamentRequest) domain.Tournament {
	rounds := s.buildRoundsFromRequests(req.Rounds)

	return domain.Tournament{
		Name:                   req.Name,
		Description:            req.Description,
		StartDate:              req.StartDate,
		EndDate:                req.EndDate,
		Status:                 domain.StatusDraft,
		Rounds:                 rounds,
		AllowUnderfilledGroups: req.AllowUnderfilledGroups,
	}
}

// buildRoundsFromRequests converts request rounds to domain rounds
func (s *TournamentService) buildRoundsFromRequests(requestRounds []requests.CreateTournamentRoundRequest) []domain.Round {
	rounds := make([]domain.Round, 0, len(requestRounds))

	for _, round := range requestRounds {
		newRound := domain.Round{
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

	return rounds
}

// handleRepositoryError handles repository errors consistently
func (s *TournamentService) handleRepositoryError(err error) {
	if err != nil {
		panic(err)
	}
}
