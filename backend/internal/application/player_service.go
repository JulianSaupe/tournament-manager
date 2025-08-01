package application

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"github.com/google/uuid"
)

type PlayerService struct {
	playerRepository output.PlayerRepository
}

func NewPlayerService(playerRepository output.PlayerRepository) input.PlayerService {
	return &PlayerService{
		playerRepository: playerRepository,
	}
}

func (s *PlayerService) CreatePlayer(name string, tournamentId string) *domain.Player {
	player := &domain.Player{
		Id:           uuid.New().String(),
		Name:         name,
		TournamentId: tournamentId,
	}

	player, err := s.playerRepository.Save(player)

	if err != nil {
		panic(err)
	}

	return player
}
