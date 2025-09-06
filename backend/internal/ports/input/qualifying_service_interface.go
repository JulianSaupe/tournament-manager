package input

import (
	"Tournament/internal/domain"
	"context"
)

type QualifyingServiceInterface interface {
	GetQualifyingByTournamentId(ctx context.Context, id string) *domain.Qualifying

	DeleteQualifyingByTournamentId(ctx context.Context, id string)
}
