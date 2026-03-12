package service

import (
	"context"
	"engine/internal/proto/authentication"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthenticationService acts as a gRPC client to an external authentication service
type AuthenticationService struct {
	client authentication.AuthenticationServiceClient
	conn   *grpc.ClientConn
}

// NewAuthenticationService creates a new authentication service client
func NewAuthenticationService(authServerAddr string) (*AuthenticationService, error) {
	conn, err := grpc.NewClient(authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to authentication service: %w", err)
	}

	client := authentication.NewAuthenticationServiceClient(conn)

	return &AuthenticationService{
		client: client,
		conn:   conn,
	}, nil
}

// ValidateSession validates a session and returns the user ID
func (s *AuthenticationService) ValidateSession(ctx context.Context, sessionID string) (string, error) {
	req := &authentication.ValidateSessionRequest{
		SessionId: sessionID,
	}

	resp, err := s.client.ValidateSession(ctx, req)
	if err != nil {
		return "", fmt.Errorf("failed to validate session: %w", err)
	}

	if !resp.Valid {
		return "", fmt.Errorf("invalid session")
	}

	return resp.UserId, nil
}

// Close closes the gRPC connection
func (s *AuthenticationService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
