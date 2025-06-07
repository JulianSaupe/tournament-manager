package http

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

// TournamentHandler handles HTTP requests for tournament operations
type TournamentHandler struct {
	tournamentService input.TournamentService
}

// NewTournamentHandler creates a new tournament handler
func NewTournamentHandler(tournamentService input.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		tournamentService: tournamentService,
	}
}

// RegisterRoutes registers the tournament routes
func (h *TournamentHandler) RegisterRoutes(router chi.Router) {
	router.Route("/tournament", func(r chi.Router) {
		r.Get("/list", h.ListTournaments)
		r.Post("/create", h.CreateTournament)
		r.Get("/{id}", h.GetTournament)
		r.Put("/{id}/status", h.UpdateTournamentStatus)
	})
}

// CreateTournamentRequest represents the request body for creating a tournament
type CreateTournamentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

// CreateTournament handles the creation of a new tournament
func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tournament, err := h.tournamentService.CreateTournament(req.Name, req.Description, req.StartDate, req.EndDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(tournament)
}

// GetTournament handles retrieving a tournament by ID
func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tournament, err := h.tournamentService.GetTournament(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tournament)
}

// ListTournaments handles retrieving all tournaments
func (h *TournamentHandler) ListTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments, err := h.tournamentService.ListTournaments()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tournaments)
}

// UpdateTournamentStatusRequest represents the request body for updating a tournament status
type UpdateTournamentStatusRequest struct {
	Status string `json:"status"`
}

// UpdateTournamentStatus handles updating the status of a tournament
func (h *TournamentHandler) UpdateTournamentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTournamentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert string status to domain.TournamentStatus
	// This is a simplification; in a real application, you would validate the status
	var status domain.TournamentStatus
	switch req.Status {
	case "DRAFT":
		status = domain.StatusDraft
	case "ACTIVE":
		status = domain.StatusActive
	case "COMPLETED":
		status = domain.StatusCompleted
	case "CANCELLED":
		status = domain.StatusCancelled
	default:
		http.Error(w, "Invalid status", http.StatusBadRequest)
		return
	}

	tournament, err := h.tournamentService.UpdateTournamentStatus(id, status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tournament)
}
