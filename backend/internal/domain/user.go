package domain

import (
	"time"
)

// User represents a user entity
type User struct {
	ID        string    `bun:"id,pk,autoincrement" json:"id"`
	Username  string    `bun:"username,unique" json:"username"`
	Password  string    `bun:"password" json:"-"` // Password is not included in JSON responses
	CreatedAt time.Time `bun:"created_at" json:"createdAt"`
	UpdatedAt time.Time `bun:"updated_at" json:"updatedAt"`
}

// UserCredentials represents the credentials used for authentication
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
