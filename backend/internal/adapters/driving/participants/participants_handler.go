package participants

import (
	"Tournament/internal/ports/input"
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
		r.Patch("/", h.UpdateParticipant)
		r.Delete("/{id}", h.DeleteParticipant)
	})
}

func (h *Handler) ListParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) GetParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) CreateParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) UpdateParticipant(w http.ResponseWriter, r *http.Request) {

}

func (h *Handler) DeleteParticipant(w http.ResponseWriter, r *http.Request) {

}
