package domain

import "github.com/uptrace/bun"

// Tournament represents a tournament entity
type Tournament struct {
	Id          string           `bun:"id,pk" json:"id"`
	Name        string           `bun:"name" json:"name"`
	Description string           `bun:"description" json:"description"`
	StartDate   string           `bun:"start_date" json:"startDate"`
	EndDate     string           `bun:"end_date" json:"endDate"`
	Status      TournamentStatus `bun:"status" json:"status"`
	Players     []Player         `bun:"rel:has-many,join:id=tournament_id" json:"players"`
	PlayerCount int              `bun:"player_count" json:"playerCount"`
	Rounds      []Round          `bun:"rel:has-many,join:id=tournament_id" json:"rounds"`
}

type Round struct {
	bun.BaseModel `bun:"table:rounds"`

	Id                     string  `bun:"id,pk" json:"id"`
	Name                   string  `bun:"name" json:"name"`
	TournamentId           string  `bun:"tournament_id" json:"tournamentId"`
	MatchCount             int     `bun:"match_count" json:"matchCount"`
	PlayerCount            int     `bun:"player_count" json:"playerCount"`
	PlayerAdvancementCount int     `bun:"player_advancement_count" json:"playerAdvancementCount"`
	GroupSize              int     `bun:"group_size" json:"groupSize"`
	Groups                 []Group `bun:"rel:has-many,join:round_id=id" json:"groups"`
}

type Group struct {
	bun.BaseModel `bun:"table:groups"`

	Id      string  `bun:"id,pk" json:"id"`
	Name    string  `bun:"name" json:"name"`
	RoundId string  `bun:"round_id" json:"roundId"`
	Matches []Match `bun:"matches,rel:has-many,join:id=group_id" json:"matches"`
}

type Match struct {
	bun.BaseModel `bun:"table:matches"`

	Id         string      `bun:"id,pk" json:"id"`
	GroupId    string      `bun:"group_id" json:"groupId"`
	MapName    string      `bun:"map_name" json:"mapName"`
	Placements []Placement `bun:"placements" json:"placements"`
}

type Placement struct {
	bun.BaseModel `bun:"table:placements"`

	Id        string `bun:"id,pk" json:"id"`
	MatchId   string `bun:"match_id" json:"matchId"`
	PlayerId  string `bun:"player_id" json:"playerId"`
	Placement int    `bun:"placement" json:"placement"`
}

type PlayerToGroup struct {
	bun.BaseModel `bun:"table:player_to_group"`

	PlayerId string  `bun:"player_id" json:"playerId"`
	Player   *Player `bun:"rel:belongs-to,join:player_id=id"`
	GroupId  string  `bun:"group_id" json:"groupId"`
	Group    *Group  `bun:"rel:belongs-to,join:group_id=id"`
}

type IndexTournament struct {
	bun.BaseModel `bun:"table:tournaments"`

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
	// StatusCancelled indicates the tournament is canceled
	StatusCancelled TournamentStatus = "CANCELLED"
)
