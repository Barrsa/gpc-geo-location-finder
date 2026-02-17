package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gpc-geo-location-finder/internal"
	"github.com/gpc-geo-location-finder/models"
)

// HandlerOptions contains options for the CheckPing handler.
type HandlerOptions struct {
	Endpoints map[string]models.Endpoint
	Client    *http.Client
	Timeout   time.Duration
}

// HandleCheckPing returns an HTTP handler function that checks all endpoints
// and returns a JSON summary of response times.
func HandleCheckPing(opts HandlerOptions) http.HandlerFunc {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	if opts.Client == nil {
		opts.Client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), opts.Timeout)
		defer cancel()

		// Add CORS headers
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-store")

		summary := internal.CheckAllEndpoints(ctx, opts.Endpoints, opts.Client)

		if err := json.NewEncoder(w).Encode(summary); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	}
}
