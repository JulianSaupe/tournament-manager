package input

import (
	"Tournament/internal/domain"
	"context"
)

type PlayerService interface {
	CreatePlayer(ctx context.Context, name string, tournamentId string) *domain.Player
}
