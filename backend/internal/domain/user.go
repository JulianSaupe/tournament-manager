package domain

import (
	"time"
)

// User represents a user entity
type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Password is not included in JSON responses
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// UserCredentials represents the credentials used for authentication
type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
