package domain

// Tournament represents a tournament entity
type Tournament struct {
	Id          string           `bun:"id,pk" json:"id"`
	Name        string           `bun:"name" json:"name"`
	Description string           `bun:"description" json:"description"`
	StartDate   string           `bun:"start_date" json:"startDate"`
	EndDate     string           `bun:"end_date" json:"endDate"`
	Status      TournamentStatus `bun:"status" json:"status"`
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
