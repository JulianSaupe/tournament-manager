package domain

type Player struct {
	// Table: players
	Id           string `json:"id"`
	Name         string `json:"name"`
	TournamentId string `json:"tournamentId"`
}
