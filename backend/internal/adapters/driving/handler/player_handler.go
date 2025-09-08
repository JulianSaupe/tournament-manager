package handler

import (
	"Tournament/internal/adapters/driving/requests"
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/adapters/driving/validation"
	"Tournament/internal/domain"
	"Tournament/internal/middleware"
	"Tournament/internal/ports/input"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type PlayerHandler struct {
	playerService     input.PlayerServiceInterface
	qualifyingService input.QualifyingServiceInterface
}

func NewPlayerHandler(
	playerService input.PlayerServiceInterface,
	qualifyingService input.QualifyingServiceInterface,
) *PlayerHandler {
	return &PlayerHandler{
		playerService:     playerService,
		qualifyingService: qualifyingService,
	}
}

func (h *PlayerHandler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.ListPlayers)
	router.Get("/{playerId}", h.GetPlayer)

	router.Group(func(router chi.Router) {
		router.Use(middleware.TournamentActiveMiddleware())
		router.Post("/", h.CreatePlayer)
		router.Patch("/{playerId}", h.UpdatePlayer)
		router.Delete("/{playerId}", h.DeletePlayer)
	})
}

func (h *PlayerHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournament := ctx.Value(middleware.TournamentKey{}).(*domain.Tournament)

	players := h.playerService.ListPlayers(ctx, tournament.Id)
	response.Send(w, r, http.StatusOK, players)
}

func (h *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "playerId")
	player := h.playerService.GetPlayer(r.Context(), id)

	response.Send(w, r, http.StatusOK, player)
}

func (h *PlayerHandler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournament := ctx.Value(middleware.TournamentKey{}).(*domain.Tournament)

	var req = validation.ValidateRequest[requests.CreatePlayerRequest](r)

	player := h.playerService.CreatePlayer(ctx, req.Name, tournament.Id)

	h.qualifyingService.AddPlayerToQualifying(ctx, tournament.Id, player.Id)

	response.Send(w, r, http.StatusCreated, player)
}

func (h *PlayerHandler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "playerId")

	var req requests.UpdatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	player := h.playerService.UpdatePlayerName(r.Context(), id, req.Name)
	response.Send(w, r, http.StatusOK, player)
}

func (h *PlayerHandler) DeletePlayer(w http.ResponseWriter, r *http.Request) {
	params, err := validation.ParseURLParams[requests.DeletePlayerRequest](r)
	if err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	h.playerService.DeletePlayer(r.Context(), params.Id)
	response.Send(w, r, http.StatusOK, nil)
}
