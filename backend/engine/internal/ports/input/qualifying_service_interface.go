package input

import (
	"context"
	"engine/internal/domain"
)

type QualifyingServiceInterface interface {
	GetQualifyingByTournamentId(ctx context.Context, id string) *domain.Qualifying

	DeleteQualifyingByTournamentId(ctx context.Context, id string)

	AddPlayerToQualifying(ctx context.Context, tournamentId string, playerId string)
}
