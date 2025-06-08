package participants

import (
	"Tournament/internal/adapters/driving/middleware"
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/domain"
	"Tournament/internal/ports/input"
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Handler struct {
	participantsService input.ParticipantsService
}

func NewParticipantsHandler(participantsService input.ParticipantsService) *Handler {
	return &Handler{
		participantsService: participantsService,
	}
}

func (h *Handler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.ListParticipant)
	router.Get("/{id}", h.GetParticipant)

	router.Group(func(router chi.Router) {
		router.Use(middleware.TournamentActiveMiddleware())
		router.Post("/", h.CreateParticipant)
		router.Patch("/{id}", h.UpdateParticipant)
		router.Delete("/{id}", h.DeleteParticipant)
	})
}

func (h *Handler) ListParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetParticipant(w http.ResponseWriter, r *http.Request) {
}

func (h *Handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {
	tournament := r.Context().Value(middleware.TournamentKey{}).(*domain.Tournament)

	var req CreateParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	participant := h.participantsService.CreateParticipant(req.Name, tournament.Id)
	response.Send(w, r, http.StatusCreated, participant)
}

func (h *Handler) UpdateParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteParticipant(w http.ResponseWriter, r *http.Request) {

}
