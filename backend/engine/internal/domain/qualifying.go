package domain

type Qualifying struct {
	TournamentId string              `json:"tournament_id"`
	Players      []*QualifyingPlayer `json:"players"`
}

type QualifyingPlayer struct {
	PlayerId   string `json:"player_id"`
	Name       string `json:"name"`
	Position   int    `json:"position"`
	SignupDate string `json:"signup_date"`
	Time       int    `json:"time"`
}
