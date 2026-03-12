package service

import (
	"context"
	"engine/internal/domain"
	"engine/internal/ports/input"
	"engine/internal/ports/output"
)

type QualifyingService struct {
	qualifyingRepository output.QualifyingRepositoryInterface
}

func NewQualifyingService(qualifyingRepository output.QualifyingRepositoryInterface) input.QualifyingServiceInterface {
	return &QualifyingService{
		qualifyingRepository: qualifyingRepository,
	}
}

func (q QualifyingService) GetQualifyingByTournamentId(ctx context.Context, id string) *domain.Qualifying {
	qualifying, err := q.qualifyingRepository.FindByTournamentId(ctx, id)
	if err != nil {
		panic(err)
	}

	return qualifying
}

func (q QualifyingService) DeleteQualifyingByTournamentId(ctx context.Context, id string) {
	panic("implement me")
}

func (q QualifyingService) AddPlayerToQualifying(ctx context.Context, tournamentId string, playerId string) {
	err := q.qualifyingRepository.AddPlayer(ctx, tournamentId, playerId)

	if err != nil {
		panic(err)
	}
}
