package postgres

import (
	"Tournament/internal/domain"
	"context"
	"errors"
	"fmt"
	"github.com/uptrace/bun"
)

type ParticipantsRepository struct {
	db *bun.DB
}

func NewParticipantsRepository(db *bun.DB) (*ParticipantsRepository, error) {
	if db == nil {
		return nil, errors.New("db cannot be nil")
	}

	db.RegisterModel((*domain.Participant)(nil))
	db.NewCreateTable().Model((*domain.Participant)(nil)).IfNotExists().Exec(context.Background())

	return &ParticipantsRepository{
		db: db,
	}, nil
}

func (r ParticipantsRepository) Save(participant *domain.Participant) (*domain.Participant, error) {
	_, err := r.db.NewInsert().Model(participant).Exec(context.Background())

	if err != nil {
		return nil, fmt.Errorf("error saving participant: %w", err)
	}

	return participant, nil
}
