package application

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"Tournament/internal/ports/output"
	"github.com/google/uuid"
)

type ParticipantsServiceImpl struct {
	participantsRepository output.ParticipantsRepository
}

func NewParticipantsService(participantsRepository output.ParticipantsRepository) input.ParticipantsService {
	return &ParticipantsServiceImpl{
		participantsRepository: participantsRepository,
	}
}

func (s *ParticipantsServiceImpl) CreateParticipant(name string, tournamentId string) *domain.Participant {
	participant := &domain.Participant{
		Id:           uuid.New().String(),
		Name:         name,
		TournamentId: tournamentId,
	}

	participant, err := s.participantsRepository.Save(participant)

	if err != nil {
		panic(err)
	}

	return participant
}
