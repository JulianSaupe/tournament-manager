package domain

// Tournament represents a tournament entity
type Tournament struct {
	// Table: tournaments
	Id                     string           `json:"id"`
	Name                   string           `json:"name"`
	Description            string           `json:"description"`
	StartDate              string           `json:"startDate"`
	EndDate                string           `json:"endDate"`
	Status                 TournamentStatus `json:"status"`
	Players                []Player         `json:"players"`
	PlayerCount            int              `json:"playerCount"`
	Rounds                 []Round          `json:"rounds"`
	AllowUnderfilledGroups bool             `json:"allowUnderfilledGroups"`
}

type Round struct {
	// Table: rounds
	Id                     string  `json:"id"`
	Name                   string  `json:"name"`
	TournamentId           string  `json:"tournamentId"`
	MatchCount             int     `json:"matchCount"`
	PlayerCount            int     `json:"playerCount"`
	PlayerAdvancementCount int     `json:"playerAdvancementCount"`
	GroupSize              int     `json:"groupSize"`
	ConcurrentGroupCount   int     `json:"concurrentGroupCount"`
	Groups                 []Group `json:"groups"`
}

type Group struct {
	// Table: groups
	Id      string  `json:"id"`
	Name    string  `json:"name"`
	RoundId string  `json:"roundId"`
	Matches []Match `json:"matches"`
}

type Match struct {
	// Table: matches
	Id         string      `json:"id"`
	GroupId    string      `json:"groupId"`
	MapName    string      `json:"mapName"`
	Placements []Placement `json:"placements"`
}

type Placement struct {
	// Table: placements
	Id        string `json:"id"`
	MatchId   string `json:"matchId"`
	PlayerId  string `json:"playerId"`
	Placement int    `json:"placement"`
}

type PlayerToGroup struct {
	// Table: player_to_group
	PlayerId string  `json:"playerId"`
	Player   *Player `json:"player,omitempty"`
	GroupId  string  `json:"groupId"`
	Group    *Group  `json:"group,omitempty"`
}

type IndexTournament struct {
	// Table: tournaments
	Id          string           `json:"id"`
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
	// StatusCancelled indicates the tournament is canceled
	StatusCancelled TournamentStatus = "CANCELLED"
)
