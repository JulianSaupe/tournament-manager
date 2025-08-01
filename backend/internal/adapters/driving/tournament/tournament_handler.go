package tournament

import (
	"Tournament/internal/adapters/driving/player"
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/domain"
	"Tournament/internal/middleware"
	"Tournament/internal/ports/input"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	tournamentService input.TournamentService
	playerService     input.PlayerService
}

func NewTournamentHandler(tournamentService input.TournamentService, playerService input.PlayerService) *Handler {
	return &Handler{
		tournamentService: tournamentService,
		playerService:     playerService,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	playerRouter := chi.NewRouter()
	playerHandler := player.NewPlayerHandler(h.playerService)
	playerHandler.RegisterRoutes(playerRouter)

	router.Route("/tournament", func(router chi.Router) {
		router.Get("/", h.ListTournaments)
		router.Post("/", h.CreateTournament)

		router.Route("/{id}", func(router chi.Router) {
			router.Use(middleware.TournamentMiddleware(h.tournamentService))
			router.Patch("/status", h.UpdateTournamentStatus)
			router.Get("/", h.GetTournament)

			router.Group(func(router chi.Router) {
				router.Use(middleware.TournamentActiveMiddleware())
				router.Delete("/", h.DeleteTournament)
			})

			router.Mount("/player", playerRouter)
		})
	})
}

func (h *Handler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	tournament := h.tournamentService.CreateTournament(ctx, req.Name, req.Description, req.StartDate, req.EndDate)

	response.Send(w, r, http.StatusCreated, tournament)
}

func (h *Handler) GetTournament(w http.ResponseWriter, r *http.Request) {
	// The tournament is already retrieved by middleware and stored in context
	tournament := r.Context().Value(middleware.TournamentKey{}).(*domain.Tournament)
	response.Send(w, r, http.StatusOK, tournament)
}

func (h *Handler) ListTournaments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournaments := h.tournamentService.ListTournaments(ctx)

	response.Send(w, r, http.StatusOK, tournaments)
}

func (h *Handler) UpdateTournamentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req UpdateTournamentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
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
		response.SendError(w, r, http.StatusBadRequest, "Invalid status")
		return
	}

	ctx := r.Context()
	tournament := h.tournamentService.UpdateTournamentStatus(ctx, id, status)

	response.Send(w, r, http.StatusOK, tournament)
}

func (h *Handler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()

	h.tournamentService.DeleteTournament(ctx, id)
}
