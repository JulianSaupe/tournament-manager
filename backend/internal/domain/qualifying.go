package domain

type Qualifying struct {
	TournamentId string              `json:"tournament_id"`
	Players      []*QualifyingPlayer `json:"players"`
}

type QualifyingPlayer struct {
	PlayerId   string `json:"player_id"`
	PlayerName string `json:"player_name"`
	Position   int    `json:"position"`
	SignupDate string `json:"signup_date"`
	BestTime   int    `json:"best_time"`
}
