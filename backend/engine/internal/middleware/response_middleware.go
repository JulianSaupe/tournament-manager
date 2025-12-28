package middleware

import (
	"Tournament/internal/adapters/driving/response"
	"Tournament/internal/domain"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"
	"runtime/debug"
)

func CustomRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				errorMessage := "Internal server error"
				var err error

				if e, ok := rvr.(error); ok {
					err = e
					errorMessage = e.Error()
				} else if str, ok := rvr.(string); ok {
					err = errors.New(str)
					errorMessage = str
				} else {
					err = errors.New(errorMessage)
				}

				statusCode := inferStatusCode(err)

				if statusCode == http.StatusInternalServerError {
					logEntry := middleware.GetLogEntry(r)
					if logEntry != nil {
						logEntry.Panic(rvr, debug.Stack())
					} else {
						_, _ = fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
						debug.PrintStack()
					}
				}

				response.SendError(w, r, statusCode, errorMessage)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// inferStatusCode determines the appropriate HTTP status code based on the error type
func inferStatusCode(err error) int {
	if domain.IsNotFound(err) {
		return http.StatusNotFound
	}
	if domain.IsInvalidParameter(err) {
		return http.StatusBadRequest
	}
	if domain.IsUnauthorized(err) {
		return http.StatusUnauthorized
	}
	if domain.IsForbidden(err) {
		return http.StatusForbidden
	}
	if domain.IsNotAllowed(err) {
		return http.StatusMethodNotAllowed
	}

	return http.StatusInternalServerError
}
