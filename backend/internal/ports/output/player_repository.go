package output

import (
	"Tournament/internal/domain"
)

type PlayerRepository interface {
	Save(player *domain.Player) (*domain.Player, error)
}
