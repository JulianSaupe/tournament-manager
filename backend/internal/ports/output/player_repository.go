package output

import (
	"Tournament/internal/domain"
	"context"
)

type PlayerRepository interface {
	Save(ctx context.Context, player *domain.Player) (*domain.Player, error)
}
