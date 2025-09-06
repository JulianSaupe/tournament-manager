package output

import (
	"Tournament/internal/domain"
	"context"
)

type QualifyingRepositoryInterface interface {
	FindByTournamentId(ctx context.Context, id string) (*domain.Qualifying, error)

	DeleteByTournamentId(ctx context.Context, id string) error
}
