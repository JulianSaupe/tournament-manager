package domain

import "github.com/uptrace/bun"

type Player struct {
	bun.BaseModel `bun:"table:players"`

	Id           string `json:"id" bun:"id,pk"`
	Name         string `json:"name" bun:"name"`
	TournamentId string `json:"tournamentId" bun:"tournament_id"`
}
