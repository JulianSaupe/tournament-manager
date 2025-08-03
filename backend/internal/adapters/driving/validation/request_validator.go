package validation

import (
	"Tournament/internal/adapters/driving/requests"
	"Tournament/internal/domain"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ValidateRequest[T any](r *http.Request) *T {
	var req T
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(domain.NewInvalidParameterError("Invalid request parameters."))
	}

	validate := validator.New()
	err := validate.Struct(&req)
	if err != nil {
		panic(domain.NewInvalidParameterError("Invalid request parameters"))
	}

	return &req
}

func ValidateCreateTournamentRequest(r *http.Request) *requests.CreateTournamentRequest {
	var req requests.CreateTournamentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(domain.NewInvalidParameterError("Invalid request parameters."))
	}

	validate := validator.New()
	err := validate.Struct(&req)
	if err != nil {
		panic(domain.NewInvalidParameterError("Invalid request parameters"))
	}

	var previousRound *requests.CreateTournamentRoundRequest

	for _, round := range req.Rounds {
		if req.AllowUnderfilledGroups == false {
			var playersInRound = round.GroupCount * round.GroupSize

			if previousRound == nil {
				if playersInRound != req.PlayerCount {
					panic(domain.NewInvalidParameterError("Number of players in first round must be equal to total players in tournament"))
				}
			} else if playersInRound != (previousRound.PlayerAdvancementCount * previousRound.GroupCount) {
				panic(domain.NewInvalidParameterError("Number of players in round must be equal to total advancing players of previous round"))
			}
		}

		if round.GroupSize <= 0 {
			panic(domain.NewInvalidParameterError("Group size must be greater than 0"))
		}

		if round.GroupCount <= 0 {
			panic(domain.NewInvalidParameterError("Group count must be greater than 0"))
		}

		if round.PlayerAdvancementCount > round.GroupSize {
			panic(domain.NewInvalidParameterError("Player advancement count cannot exceed total players in group"))
		}

		previousRound = &round
	}

	return &req
}
