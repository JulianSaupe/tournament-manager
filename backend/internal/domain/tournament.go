package domain

// Tournament represents a tournament entity
type Tournament struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	StartDate   string           `json:"startDate"`
	EndDate     string           `json:"endDate"`
	Status      TournamentStatus `json:"status"`
}

// TournamentStatus represents the status of a tournament
type TournamentStatus string

const (
	// StatusDraft indicates the tournament is in draft mode
	StatusDraft TournamentStatus = "DRAFT"
	// StatusActive indicates the tournament is active
	StatusActive TournamentStatus = "ACTIVE"
	// StatusCompleted indicates the tournament is completed
	StatusCompleted TournamentStatus = "COMPLETED"
	// StatusCancelled indicates the tournament is cancelled
	StatusCancelled TournamentStatus = "CANCELLED"
)
