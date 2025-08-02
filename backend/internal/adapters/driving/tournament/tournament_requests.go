package tournament

type CreateTournamentRequest struct {
	Name        string                         `json:"name" validate:"required,min=3,max=255"`
	Description string                         `json:"description" validate:"required,min=3,max=255"`
	StartDate   string                         `json:"startDate" validate:"required"`
	EndDate     string                         `json:"endDate" validate:"required"`
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
	Status string `json:"status" validate:"required"`
}
