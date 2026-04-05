package service

import (
	"context"
	"engine/internal/proto/authorization"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// AuthorizationService acts as a gRPC client to an external authorization service
type AuthorizationService struct {
	client authorization.AuthorizationServiceClient
	conn   *grpc.ClientConn
}

// NewAuthorizationService creates a new authorization service client
func NewAuthorizationService(authServerAddr string) (*AuthorizationService, error) {
	// Connect to the authorization gRPC server
	conn, err := grpc.NewClient(authServerAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to authorization service: %w", err)
	}

	client := authorization.NewAuthorizationServiceClient(conn)

	return &AuthorizationService{
		client: client,
		conn:   conn,
	}, nil
}

// CheckPermission checks if a user has permission for a specific resource and action
func (s *AuthorizationService) CheckPermission(ctx context.Context, userID, name string) (bool, string, error) {
	req := &authorization.CheckPermissionRequest{
		UserId:   userID,
		PermissionName: name,
	}

	resp, err := s.client.CheckPermission(ctx, req)
	if err != nil {
		return false, "", fmt.Errorf("failed to check permission: %w", err)
	}

	return resp.Allowed, resp.Message, nil
}

// Close closes the gRPC connection
func (s *AuthorizationService) Close() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
