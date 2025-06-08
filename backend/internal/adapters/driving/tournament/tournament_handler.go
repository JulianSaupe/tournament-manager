package tournament

import (
	. "Tournament/internal/adapters/driving/response"
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	tournamentService input.TournamentService
}

func NewTournamentHandler(tournamentService input.TournamentService) *Handler {
	return &Handler{
		tournamentService: tournamentService,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Route("/tournament", func(r chi.Router) {
		r.Get("/list", h.ListTournaments)
		r.Post("/create", h.CreateTournament)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetTournament)
			r.Patch("/status", h.UpdateTournamentStatus)
			r.Delete("/", h.DeleteTournament)

			r.Route("/participants", func(r chi.Router) {
				participantsRouter := chi.NewRouter()
				r.Mount("/", participantsRouter)

			})
		})
	})
}

func (h *Handler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	tournament := h.tournamentService.CreateTournament(req.Name, req.Description, req.StartDate, req.EndDate)

	Send(w, r, http.StatusCreated, tournament)
}

func (h *Handler) GetTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	tournament := h.tournamentService.GetTournament(id)

	Send(w, r, http.StatusOK, tournament)
}

func (h *Handler) ListTournaments(w http.ResponseWriter, r *http.Request) {
	tournaments := h.tournamentService.ListTournaments()

	Send(w, r, http.StatusOK, tournaments)
}

func (h *Handler) UpdateTournamentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTournamentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendError(w, r, http.StatusBadRequest, err.Error())
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
		SendError(w, r, http.StatusBadRequest, "Invalid status")
		return
	}

	tournament := h.tournamentService.UpdateTournamentStatus(id, status)

	Send(w, r, http.StatusOK, tournament)
}

func (h *Handler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	h.tournamentService.DeleteTournament(id)
}
