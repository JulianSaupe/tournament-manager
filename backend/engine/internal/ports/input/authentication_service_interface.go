package input

import "context"

// AuthenticationServiceInterface defines the contract for authentication operations
type AuthenticationServiceInterface interface {
	ValidateSession(ctx context.Context, sessionID string) (string, error)
}
