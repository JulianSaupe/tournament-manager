package http

import (
	"context"
	"encoding/json"
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

func CustomRecoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				logEntry := middleware.GetLogEntry(r)
				if logEntry != nil {
					logEntry.Panic(rvr, debug.Stack())
				} else {
					fmt.Fprintf(os.Stderr, "Panic: %+v\n", rvr)
					debug.PrintStack()
				}

				errorMessage := "Internal server error"
				if err, ok := rvr.(error); ok {
					errorMessage = err.Error()
				} else if str, ok := rvr.(string); ok {
					errorMessage = str
				}

				SendErrorResponse(w, r, http.StatusInternalServerError, errorMessage)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
