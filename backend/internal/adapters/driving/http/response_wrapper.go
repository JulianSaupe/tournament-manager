package http

import (
	"Tournament/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type StandardResponse struct {
	Server     string      `json:"server"`
	StartTime  time.Time   `json:"startTime"`
	EndTime    time.Time   `json:"endTime"`
	StatusCode int         `json:"statusCode"`
	Data       interface{} `json:"data"`
}

type requestStartTimeKey struct{}

func RequestStartTimeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), requestStartTimeKey{}, time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRequestStartTime(r *http.Request) time.Time {
	if startTime, ok := r.Context().Value(requestStartTimeKey{}).(time.Time); ok {
		return startTime
	}
	return time.Now()
}

func NewResponse(r *http.Request, status int, data interface{}) StandardResponse {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}

	requestID := middleware.GetReqID(r.Context())
	if requestID != "" {
		hostname = hostname + "-" + requestID
	}

	return StandardResponse{
		Server:     hostname,
		StartTime:  getRequestStartTime(r),
		EndTime:    time.Now(),
		StatusCode: status,
		Data:       data,
	}
}

func SendResponse(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	response := NewResponse(r, status, data)
	response.EndTime = time.Now()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		return
	}
}

func SendErrorResponse(w http.ResponseWriter, r *http.Request, status int, message string) {
	errorData := map[string]string{"error": message}
	SendResponse(w, r, status, errorData)
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

	return http.StatusInternalServerError
}

func CustomRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					_, _ = fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}

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
				SendErrorResponse(w, r, statusCode, errorMessage)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
