package validation

import (
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
