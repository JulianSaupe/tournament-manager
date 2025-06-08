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
	router.Route("/participants", func(r chi.Router) {
		r.Get("/", h.ListParticipant)
		r.Get("/{id}", h.GetParticipant)
		r.Post("/", h.CreateParticipant)
		r.Patch("/{id}", h.UpdateParticipant)
		r.Delete("/{id}", h.DeleteParticipant)
	})
}

func (h *Handler) ListParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetParticipant(w http.ResponseWriter, r *http.Request) {
	tournament := r.Context().Value(middleware.TournamentKey{}).(*domain.Tournament)

	var req CreateParticipantRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.SendError(w, r, http.StatusBadRequest, err.Error())
		return
	}

	h.participantsService.CreateParticipant(req.Name, tournament.Id)
}

func (h *Handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteParticipant(w http.ResponseWriter, r *http.Request) {

}
