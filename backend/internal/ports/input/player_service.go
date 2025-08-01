package input

import "Tournament/internal/domain"

type PlayerService interface {
	CreatePlayer(name string, tournamentId string) *domain.Player
}
