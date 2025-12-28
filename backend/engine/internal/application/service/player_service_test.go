package service

import (
	"Tournament/internal/domain"
	"context"
	"errors"
	"testing"
)

// MockPlayerRepository is a mock implementation of the PlayerRepositoryInterface
type MockPlayerRepository struct {
	players map[string]*domain.Player
	// Track method calls for verification
	insertCalled     bool
	deleteCalled     bool
	findAllCalled    bool
	findByIDCalled   bool
	updateNameCalled bool
	// Control error responses for testing error paths
	shouldReturnError bool
}

// NewMockPlayerRepository creates a new mock repository with optional initial players
func NewMockPlayerRepository(initialPlayers []*domain.Player, shouldReturnError bool) *MockPlayerRepository {
	playersMap := make(map[string]*domain.Player)
	for _, player := range initialPlayers {
		playersMap[player.Id] = player
	}

	return &MockPlayerRepository{
		players:           playersMap,
		shouldReturnError: shouldReturnError,
	}
}

// InsertNewPlayer mocks inserting a new player
func (m *MockPlayerRepository) InsertNewPlayer(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	m.insertCalled = true

	if m.shouldReturnError {
		return nil, errors.New("mock insert error")
	}

	// Simulate ID generation if not provided
	if player.Id == "" {
		player.Id = "mock-id-" + player.Name
	}

	m.players[player.Id] = player
	return player, nil
}

// Delete mocks deleting a player
func (m *MockPlayerRepository) Delete(ctx context.Context, id string) error {
	m.deleteCalled = true

	if m.shouldReturnError {
		return errors.New("mock delete error")
	}

	delete(m.players, id)
	return nil
}

// FindAll mocks finding all players for a tournament
func (m *MockPlayerRepository) FindAll(ctx context.Context, tournamentId string) ([]*domain.Player, error) {
	m.findAllCalled = true

	if m.shouldReturnError {
		return nil, errors.New("mock find all error")
	}

	var result []*domain.Player
	for _, player := range m.players {
		if player.TournamentId == tournamentId {
			result = append(result, player)
		}
	}

	return result, nil
}

// FindByID mocks finding a player by ID
func (m *MockPlayerRepository) FindByID(ctx context.Context, id string) (*domain.Player, error) {
	m.findByIDCalled = true

	if m.shouldReturnError {
		return nil, errors.New("mock find by ID error")
	}

	player, exists := m.players[id]
	if !exists {
		return nil, errors.New("player not found")
	}

	return player, nil
}

// UpdateName mocks updating a player's name
func (m *MockPlayerRepository) UpdateName(ctx context.Context, player *domain.Player) (*domain.Player, error) {
	m.updateNameCalled = true

	if m.shouldReturnError {
		return nil, errors.New("mock update error")
	}

	existingPlayer, exists := m.players[player.Id]
	if !exists {
		return nil, errors.New("player not found")
	}

	existingPlayer.Name = player.Name
	return existingPlayer, nil
}

// Test cases for PlayerService

func TestCreatePlayer(t *testing.T) {
	// Test successful creation
	t.Run("successful creation", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockPlayerRepository(nil, false)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act
		player := service.CreatePlayer(ctx, "Test Player", "tournament-123")

		// Assert
		if player == nil {
			t.Fatal("Expected player to be created, got nil")
		}
		if player.Name != "Test Player" {
			t.Errorf("Expected player name to be 'Test Player', got '%s'", player.Name)
		}
		if player.TournamentId != "tournament-123" {
			t.Errorf("Expected tournament ID to be 'tournament-123', got '%s'", player.TournamentId)
		}
		if !mockRepo.insertCalled {
			t.Error("Expected InsertNewPlayer to be called")
		}
	})

	// Test error handling
	t.Run("repository error", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockPlayerRepository(nil, true)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act & Assert
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected function to panic, but it didn't")
			}
		}()

		service.CreatePlayer(ctx, "Test Player", "tournament-123")
	})
}

func TestDeletePlayer(t *testing.T) {
	// Test successful deletion
	t.Run("successful deletion", func(t *testing.T) {
		// Arrange
		initialPlayers := []*domain.Player{
			{Id: "player-123", Name: "Test Player", TournamentId: "tournament-123"},
		}
		mockRepo := NewMockPlayerRepository(initialPlayers, false)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act
		service.DeletePlayer(ctx, "player-123")

		// Assert
		if !mockRepo.deleteCalled {
			t.Error("Expected Delete to be called")
		}
	})

	// Test error handling
	t.Run("repository error", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockPlayerRepository(nil, true)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act & Assert
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected function to panic, but it didn't")
			}
		}()

		service.DeletePlayer(ctx, "player-123")
	})
}

func TestListPlayers(t *testing.T) {
	// Test successful listing
	t.Run("successful listing", func(t *testing.T) {
		// Arrange
		initialPlayers := []*domain.Player{
			{Id: "player-1", Name: "Player 1", TournamentId: "tournament-123"},
			{Id: "player-2", Name: "Player 2", TournamentId: "tournament-123"},
			{Id: "player-3", Name: "Player 3", TournamentId: "tournament-456"},
		}
		mockRepo := NewMockPlayerRepository(initialPlayers, false)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act
		players := service.ListPlayers(ctx, "tournament-123")

		// Assert
		if !mockRepo.findAllCalled {
			t.Error("Expected FindAll to be called")
		}
		if len(players) != 2 {
			t.Errorf("Expected 2 players, got %d", len(players))
		}
	})

	// Test error handling
	t.Run("repository error", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockPlayerRepository(nil, true)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act & Assert
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected function to panic, but it didn't")
			}
		}()

		service.ListPlayers(ctx, "tournament-123")
	})
}

func TestGetPlayer(t *testing.T) {
	// Test successful retrieval
	t.Run("successful retrieval", func(t *testing.T) {
		// Arrange
		initialPlayers := []*domain.Player{
			{Id: "player-123", Name: "Test Player", TournamentId: "tournament-123"},
		}
		mockRepo := NewMockPlayerRepository(initialPlayers, false)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act
		player := service.GetPlayer(ctx, "player-123")

		// Assert
		if !mockRepo.findByIDCalled {
			t.Error("Expected FindByID to be called")
		}
		if player == nil {
			t.Fatal("Expected player to be found, got nil")
		}
		if player.Id != "player-123" {
			t.Errorf("Expected player ID to be 'player-123', got '%s'", player.Id)
		}
		if player.Name != "Test Player" {
			t.Errorf("Expected player name to be 'Test Player', got '%s'", player.Name)
		}
	})

	// Test error handling
	t.Run("repository error", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockPlayerRepository(nil, true)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act & Assert
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected function to panic, but it didn't")
			}
		}()

		service.GetPlayer(ctx, "player-123")
	})
}

func TestUpdatePlayerName(t *testing.T) {
	// Test successful update
	t.Run("successful update", func(t *testing.T) {
		// Arrange
		initialPlayers := []*domain.Player{
			{Id: "player-123", Name: "Old Name", TournamentId: "tournament-123"},
		}
		mockRepo := NewMockPlayerRepository(initialPlayers, false)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act
		player := service.UpdatePlayerName(ctx, "player-123", "New Name")

		// Assert
		if !mockRepo.updateNameCalled {
			t.Error("Expected UpdateName to be called")
		}
		if player == nil {
			t.Fatal("Expected player to be updated, got nil")
		}
		if player.Name != "New Name" {
			t.Errorf("Expected player name to be 'New Name', got '%s'", player.Name)
		}
	})

	// Test error handling
	t.Run("repository error", func(t *testing.T) {
		// Arrange
		mockRepo := NewMockPlayerRepository(nil, true)
		service := NewPlayerService(mockRepo)
		ctx := context.Background()

		// Act & Assert
		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected function to panic, but it didn't")
			}
		}()

		service.UpdatePlayerName(ctx, "player-123", "New Name")
	})
}
