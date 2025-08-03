package handler

import (
	"Tournament/internal/adapters/driving/requests"
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/adapters/driving/validation"
	"Tournament/internal/domain"
	"Tournament/internal/middleware"
	"Tournament/internal/ports/input"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type TournamentHandler struct {
	tournamentService input.TournamentServiceInterface
	playerService     input.PlayerServiceInterface
}

func NewTournamentHandler(tournamentService input.TournamentServiceInterface, playerService input.PlayerServiceInterface) *TournamentHandler {
	return &TournamentHandler{
		tournamentService: tournamentService,
		playerService:     playerService,
	}
}

func (h *TournamentHandler) RegisterRoutes(router chi.Router) {
	playerRouter := chi.NewRouter()
	playerHandler := NewPlayerHandler(h.playerService)
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

func (h *TournamentHandler) CreateTournament(w http.ResponseWriter, r *http.Request) {
	var req = validation.ValidateCreateTournamentRequest(r)
	ctx := r.Context()
	tournament := h.tournamentService.CreateTournament(ctx, req)
	response.Send(w, r, http.StatusCreated, tournament)
}

func (h *TournamentHandler) GetTournament(w http.ResponseWriter, r *http.Request) {
	// The tournament is already retrieved by middleware and stored in context
	tournament := r.Context().Value(middleware.TournamentKey{}).(*domain.Tournament)
	response.Send(w, r, http.StatusOK, tournament)
}

func (h *TournamentHandler) ListTournaments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournaments := h.tournamentService.ListTournaments(ctx)
	response.Send(w, r, http.StatusOK, tournaments)
}

func (h *TournamentHandler) UpdateTournamentStatus(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req requests.UpdateTournamentStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	status, err := h.parseStatusString(req.Status)
	if err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	ctx := r.Context()
	tournament := h.tournamentService.UpdateTournamentStatus(ctx, id, status)
	response.Send(w, r, http.StatusOK, tournament)
}

func (h *TournamentHandler) DeleteTournament(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	h.tournamentService.DeleteTournament(ctx, id)
	response.Send(w, r, http.StatusOK, nil)
}

func (h *TournamentHandler) parseStatusString(statusStr string) (domain.TournamentStatus, error) {
	switch statusStr {
	case "DRAFT":
		return domain.StatusDraft, nil
	case "ACTIVE":
		return domain.StatusActive, nil
	case "COMPLETED":
		return domain.StatusCompleted, nil
	case "CANCELLED":
		return domain.StatusCancelled, nil
	default:
		return "", errors.New("invalid status")
	}
}
