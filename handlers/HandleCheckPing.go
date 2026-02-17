package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gpc-geo-location-finder/internal"
	"github.com/gpc-geo-location-finder/models"
)

// HandlerOptions contains options for the CheckPing handler.
type HandlerOptions struct {
	Endpoints map[string]models.Endpoint
	Client    *http.Client
	Timeout   time.Duration
}

// HandleCheckPing returns a Gin handler function that checks all endpoints
// and returns a JSON summary of response times.
func HandleCheckPing(opts HandlerOptions) gin.HandlerFunc {
	if opts.Timeout == 0 {
		opts.Timeout = 30 * time.Second
	}

	if opts.Client == nil {
		opts.Client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), opts.Timeout)
		defer cancel()

		// Add CORS headers
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Cache-Control", "no-store")

		summary := internal.CheckAllEndpoints(ctx, opts.Endpoints, opts.Client)

		c.JSON(http.StatusOK, summary)
	}
}
