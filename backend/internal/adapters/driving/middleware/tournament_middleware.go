package middleware

import (
	"Tournament/internal/ports/input"
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type TournamentKey struct{}

func TournamentMiddleware(tournamentService input.TournamentService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			tournament := tournamentService.GetTournament(id)

			ctx := context.WithValue(r.Context(), TournamentKey{}, tournament)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
