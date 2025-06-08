package http

import (
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type TournamentHandler struct {
	tournamentService input.TournamentService
}

func NewTournamentHandler(tournamentService input.TournamentService) *TournamentHandler {
	return &TournamentHandler{
		tournamentService: tournamentService,
	}
}

func (h *TournamentHandler) RegisterRoutes(router chi.Router) {
	router.Route("/tournament", func(r chi.Router) {
		r.Get("/list", h.ListTournaments)
		r.Post("/create", h.CreateTournament)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetTournament)
			r.Patch("/status", h.UpdateTournamentStatus)
			r.Delete("/", h.DeleteTournament)
		})
	})
}

type CreateTournamentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}

func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tournament := h.tournamentService.CreateTournament(req.Name, req.Description, req.StartDate, req.EndDate)

	SendResponse(w, r, http.StatusCreated, tournament)
}

func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tournament := h.tournamentService.GetTournament(id)

	SendResponse(w, r, http.StatusOK, tournament)
}

func (h *TournamentHandler) ListTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments := h.tournamentService.ListTournaments()

	SendResponse(w, r, http.StatusOK, tournaments)
}

type UpdateTournamentStatusRequest struct {
	Status string `json:"status"`
}

func (h *TournamentHandler) UpdateTournamentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTournamentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendErrorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

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
		SendErrorResponse(w, r, http.StatusBadRequest, "Invalid status")
		return
	}

	tournament := h.tournamentService.UpdateTournamentStatus(id, status)

	SendResponse(w, r, http.StatusOK, tournament)
}

func (h *TournamentHandler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.tournamentService.DeleteTournament(id)
}
