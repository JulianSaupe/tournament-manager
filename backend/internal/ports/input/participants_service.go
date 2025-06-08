package input

import "Tournament/internal/domain"

type ParticipantsService interface {
	CreateParticipant(name string, tournamentId string) *domain.Participant
}
