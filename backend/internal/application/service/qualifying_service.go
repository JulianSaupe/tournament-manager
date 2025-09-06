package service

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
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
