package player

import (
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/domain"
	"Tournament/internal/middleware"
	"Tournament/internal/ports/input"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	playerService input.PlayerService
}

func NewPlayerHandler(playerService input.PlayerService) *Handler {
	return &Handler{
		playerService: playerService,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.ListPlayers)
	router.Get("/{id}", h.GetPlayer)

	router.Group(func(router chi.Router) {
		router.Use(middleware.TournamentActiveMiddleware())
		router.Post("/", h.CreatePlayer)
		router.Patch("/{id}", h.UpdatePlayer)
		router.Delete("/{id}", h.DeletePlayer)
	})
}

func (h *Handler) ListPlayers(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tournament := ctx.Value(middleware.TournamentKey{}).(*domain.Tournament)

	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	player := h.playerService.CreatePlayer(ctx, req.Name, tournament.Id)
	response.Send(w, r, http.StatusCreated, player)
}

func (h *Handler) UpdatePlayer(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeletePlayer(w http.ResponseWriter, r *http.Request) {

}
