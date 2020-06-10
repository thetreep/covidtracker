package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/thetreep/covidtracker"
	"github.com/thetreep/covidtracker/logger"
)

// errorResponse is a generic response for sending a error.
type errorResponse struct {
	Err string `json:"err,omitempty"`
}

// Error writes an API error message to the response and logger.
func Error(ctx context.Context, w http.ResponseWriter, err error, code int) {
	// Log error.
	logger.DefaultLogger.Info(ctx, "http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal.
	if code == http.StatusInternalServerError {
		err = covidtracker.ErrInternal
	}

	// Write generic error response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}
