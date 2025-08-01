package output

import (
	"Tournament/internal/domain"
	"context"
)

type PlayerRepository interface {
	Save(ctx context.Context, player *domain.Player) (*domain.Player, error)

	Delete(ctx context.Context, id string) error

	FindAll(ctx context.Context, tournamentId string) ([]*domain.Player, error)

	FindByID(ctx context.Context, id string) (*domain.Player, error)

	UpdateName(ctx context.Context, player *domain.Player) (*domain.Player, error)
}
