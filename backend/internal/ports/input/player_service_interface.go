package input

import (
	"Tournament/internal/domain"
	"context"
)

type PlayerServiceInterface interface {
	CreatePlayer(ctx context.Context, name string, tournamentId string) *domain.Player

	DeletePlayer(ctx context.Context, id string)

	ListPlayers(ctx context.Context, tournamentId string) []*domain.Player

	GetPlayer(ctx context.Context, id string) *domain.Player

	UpdatePlayerName(ctx context.Context, id string, name string) *domain.Player
}
