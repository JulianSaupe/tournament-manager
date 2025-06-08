package domain

type Participant struct {
	Id           string `json:"id" bun:"id,pk"`
	Name         string `json:"name" bun:"name"`
	TournamentId string `json:"tournamentId" bun:"tournament_id"`
}
