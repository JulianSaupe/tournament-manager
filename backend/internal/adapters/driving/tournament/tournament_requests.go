package tournament

type CreateTournamentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

type UpdateTournamentStatusRequest struct {
	Status string `json:"status"`
}
