package service

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"context"
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

func (s *PlayerService) CreatePlayer(ctx context.Context, name string, tournamentId string) *domain.Player {
	player := &domain.Player{
		Id:           uuid.New().String(),
		Name:         name,
		TournamentId: tournamentId,
	}

	player, err := s.playerRepository.Save(ctx, player)

	if err != nil {
		panic(err)
	}

	return player
}

func (s *PlayerService) DeletePlayer(ctx context.Context, id string) {
	err := s.playerRepository.Delete(ctx, id)

	if err != nil {
		panic(err)
	}
}

func (s *PlayerService) ListPlayers(ctx context.Context, tournamentId string) []*domain.Player {
	players, err := s.playerRepository.FindAll(ctx, tournamentId)

	if err != nil {
		panic(err)
	}

	return players
}

func (s *PlayerService) GetPlayer(ctx context.Context, id string) *domain.Player {
	player, err := s.playerRepository.FindByID(ctx, id)

	if err != nil {
		panic(err)
	}

	return player
}

func (s *PlayerService) UpdatePlayerName(ctx context.Context, id string, name string) *domain.Player {
	player, err := s.playerRepository.UpdateName(ctx, &domain.Player{Id: id, Name: name})

	if err != nil {
		panic(err)
	}

	return player
}
