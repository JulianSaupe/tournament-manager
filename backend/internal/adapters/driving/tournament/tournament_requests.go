package tournament

type CreateTournamentRequest struct {
	Name        string                         `json:"name"`
	Description string                         `json:"description"`
	StartDate   string                         `json:"startDate"`
	EndDate     string                         `json:"endDate"`
	Rounds      []CreateTournamentRoundRequest `json:"rounds"`
}

type CreateTournamentRoundRequest struct {
	Name                   string `json:"name"`
	MatchCount             int    `json:"matchCount"`
	PlayerAdvancementCount int    `json:"playerAdvancementCount"`
	PlayerCount            int    `json:"playerCount"`
	GroupSize              int    `json:"groupSize"`
}

type UpdateTournamentStatusRequest struct {
	Status string `json:"status"`
}
