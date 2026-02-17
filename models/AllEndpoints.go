package models

import (
	"embed"
	"encoding/json"
	"os"
	"sync"
)

// Endpoint represents a Cloud Run service deploy in a particular region.
type Endpoint struct {
	// URL is the HTTPS URL of the service
	URL string
	// Region is the programmatic name of the region where the endpoint is
	// deployed, e.g., us-central1.
	Region string
	// RegionName is the geographic name of the region, e.g., Iowa.
	RegionName string
}

var endpointsJSON embed.FS

var (
	allEndpointsOnce sync.Once
	allEndpoints     map[string]Endpoint
)

// AllEndpoints returns the map of all endpoints, loading from JSON file if available.
func AllEndpoints() map[string]Endpoint {
	allEndpointsOnce.Do(func() {
		allEndpoints = loadEndpointsFromJSON()
	})
	return allEndpoints
}

// loadEndpointsFromJSON loads endpoints from the file specified by ENDPOINTS_FILE_PATH environment variable.
func loadEndpointsFromJSON() map[string]Endpoint {
	// Load path from environment variable
	envPath := os.Getenv("ENDPOINTS_FILE_PATH")
	if envPath == "" {
		return nil
	}

	// Try to load from the file specified in environment variable
	if data, err := os.ReadFile(envPath); err == nil {
		var endpoints map[string]Endpoint
		if err := json.Unmarshal(data, &endpoints); err == nil {
			return endpoints
		}
	}

	return nil
}
