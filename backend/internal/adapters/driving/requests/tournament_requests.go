package requests

type CreateTournamentRequest struct {
	Name                   string                         `json:"name" validate:"required,min=3,max=255"`
	Description            string                         `json:"description" validate:"required,min=3,max=255"`
	StartDate              string                         `json:"startDate" validate:"required"`
	EndDate                string                         `json:"endDate" validate:"required"`
	AllowUnderfilledGroups bool                           `json:"allowUnderfilledGroups"`
	PlayerCount            int                            `json:"playerCount" validate:"required,min=1"`
	Rounds                 []CreateTournamentRoundRequest `json:"rounds"`
}

type CreateTournamentRoundRequest struct {
	Name                   string `json:"name" validate:"required,min=3,max=255"`
	MatchCount             int    `json:"matchCount" validate:"required,min=1"`
	PlayerAdvancementCount int    `json:"playerAdvancementCount" validate:"required,min=0"`
	GroupSize              int    `json:"groupSize" validate:"required,min=2"`
	GroupCount             int    `json:"groupCount" validate:"required,min=1"`
	ConcurrentGroupCount   int    `json:"concurrentGroupCount" validate:"required,min=1"`
}

type UpdateTournamentStatusRequest struct {
	Status string `json:"status" validate:"required"`
}
