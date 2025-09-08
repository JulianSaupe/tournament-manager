package handler

import (
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/ports/input"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type QualifyingHandler struct {
	qualifyingService input.QualifyingServiceInterface
}

func NewQualifyingHandler(qualifyingService input.QualifyingServiceInterface) *QualifyingHandler {
	return &QualifyingHandler{
		qualifyingService: qualifyingService,
	}
}

func (h *QualifyingHandler) RegisterRoutes(router chi.Router) {
	router.Get("/", h.GetQualifying)
}

func (h *QualifyingHandler) GetQualifying(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ctx := r.Context()
	qualifying := h.qualifyingService.GetQualifyingByTournamentId(ctx, id)
	response.Send(w, r, http.StatusOK, qualifying)
}
