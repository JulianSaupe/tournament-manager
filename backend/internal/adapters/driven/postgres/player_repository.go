package postgres

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/output"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

type PlayerRepository struct {
	db *bun.DB
}

func NewPlayerRepository(db *bun.DB) (output.PlayerRepository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	db.RegisterModel((*domain.Player)(nil))
	db.NewCreateTable().Model((*domain.Player)(nil)).IfNotExists().Exec(context.Background())

	return &PlayerRepository{
		db: db,
	}, nil
}

func (r *PlayerRepository) Save(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	_, err := r.db.NewInsert().Model(player).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("error saving player: %w", err)
	}

	return player, nil
}

func (r *PlayerRepository) Delete(ctx context.Context, id string) error {
	result, err := r.db.NewDelete().
		Model((*domain.Player)(nil)).
		Where("id = ?", id).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("error deleting player: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.NewNotFoundError("player not found")
	}

	return nil
}

func (r *PlayerRepository) FindAll(ctx context.Context, tournamentId string) ([]*domain.Player, error) {
	players := make([]*domain.Player, 0)

	err := r.db.NewSelect().
		Model(&players).
		Where("tournament_id = ?", tournamentId).
		Scan(ctx)

	if err != nil {
		return nil, fmt.Errorf("error querying players: %w", err)
	}

	return players, nil
}

func (r *PlayerRepository) FindByID(ctx context.Context, id string) (*domain.Player, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	player := new(domain.Player)

	err := r.db.NewSelect().
		Model(player).
		Where("id = ?", id).
		Scan(ctx)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("player not found")
		}
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, fmt.Errorf("database query timed out: %w", err)
		}
		return nil, fmt.Errorf("error finding player: %w", err)
	}

	return player, nil
}

func (r *PlayerRepository) UpdateName(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	err := r.db.RunInTx(ctx, &sql.TxOptions{}, func(ctx context.Context, tx bun.Tx) error {
		result, err := tx.NewUpdate().
			Model(player).
			Set("name = ?", player.Name).
			WherePK().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("error updating player: %w", err)
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("error getting rows affected: %w", err)
		}

		if rowsAffected == 0 {
			return domain.NewNotFoundError("player not found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return player, nil
}
