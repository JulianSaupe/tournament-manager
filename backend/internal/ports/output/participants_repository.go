package output

import (
	"Tournament/internal/domain"
)

type ParticipantsRepository interface {
	Save(participant *domain.Participant) (*domain.Participant, error)
}
