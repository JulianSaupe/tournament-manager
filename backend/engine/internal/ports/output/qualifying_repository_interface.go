package output

import (
	"context"
	"engine/internal/domain"
)

type QualifyingRepositoryInterface interface {
	FindByTournamentId(ctx context.Context, id string) (*domain.Qualifying, error)

	DeleteByTournamentId(ctx context.Context, id string) error

	AddPlayer(ctx context.Context, tournamentId string, playerId string) error
}
