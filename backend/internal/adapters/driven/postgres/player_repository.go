package postgres

import (
	"Tournament/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type PlayerRepository struct {
	db *bun.DB
}

func NewPlayerRepository(db *bun.DB) (*PlayerRepository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	db.RegisterModel((*domain.Player)(nil))
	db.NewCreateTable().Model((*domain.Player)(nil)).IfNotExists().Exec(context.Background())

	return &PlayerRepository{
		db: db,
	}, nil
}

func (r PlayerRepository) Save(player *domain.Player) (*domain.Player, error) {
	_, err := r.db.NewInsert().Model(player).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error saving player: %w", err)
	}

	return player, nil
}
