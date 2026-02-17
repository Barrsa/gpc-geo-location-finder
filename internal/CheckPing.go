package internal

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gpc-geo-location-finder/models"
)

// PingResult represents the result of pinging a single endpoint.
type PingResult struct {
	Region     string `json:"region"`
	RegionName string `json:"regionName"`
	URL        string `json:"url"`
	LatencyMS  int64  `json:"latencyMs"`
	Success    bool   `json:"success"`
	Error      string `json:"error,omitempty"`
}

// SummaryResponse represents the JSON response with all ping results.
type SummaryResponse struct {
	Timestamp    string       `json:"timestamp"`
	TotalRegions int          `json:"totalRegions"`
	Successful   int          `json:"successful"`
	Failed       int          `json:"failed"`
	Fastest      *PingResult  `json:"fastest,omitempty"`
	Slowest      *PingResult  `json:"slowest,omitempty"`
	Results      []PingResult `json:"results"`
}

// CheckAllEndpoints pings all endpoints from the provided endpoints map
// and returns a summary response with response times.
func CheckAllEndpoints(ctx context.Context, endpoints map[string]models.Endpoint, client *http.Client) *SummaryResponse {
	if client == nil {
		client = &http.Client{
			Timeout: 10 * time.Second,
		}
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]PingResult, 0, len(endpoints))

	// Ping all endpoints concurrently
	for key, endpoint := range endpoints {
		wg.Add(1)
		go func(regionKey string, ep models.Endpoint) {
			defer wg.Done()

			result := pingEndpoint(ctx, client, regionKey, ep)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}(key, endpoint)
	}

	wg.Wait()

	// Build summary
	summary := buildSummary(results)
	return summary
}

// pingEndpoint pings a single endpoint and returns the result.
func pingEndpoint(ctx context.Context, client *http.Client, regionKey string, endpoint models.Endpoint) PingResult {
	result := PingResult{
		Region:     endpoint.Region,
		RegionName: endpoint.RegionName,
		URL:        endpoint.URL,
		Success:    false,
	}

	pingURL := endpoint.URL + "/api/ping"
	if endpoint.URL[len(endpoint.URL)-1] == '/' {
		pingURL = endpoint.URL + "api/ping"
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pingURL, nil)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		result.Error = err.Error()
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Error = "status code: " + http.StatusText(resp.StatusCode)
		return result
	}

	// Read the response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Error = "failed to read response: " + err.Error()
		return result
	}

	bodyStr := strings.TrimSpace(string(bodyBytes))

	// Try to parse as JSON first
	var pingResp struct {
		Region string `json:"region"`
	}
	if err := json.Unmarshal(bodyBytes, &pingResp); err == nil {
		// Successfully parsed as JSON
		result.LatencyMS = duration.Milliseconds()
		result.Success = true
		return result
	}

	// If JSON parsing fails, treat as plain text (backward compatibility)
	if bodyStr != "" {
		result.LatencyMS = duration.Milliseconds()
		result.Success = true
		return result
	}

	// If we get here, response was empty or invalid
	result.Error = "empty or invalid response"
	return result
}

// buildSummary creates a summary response from ping results.
func buildSummary(results []PingResult) *SummaryResponse {
	summary := &SummaryResponse{
		Timestamp:    time.Now().UTC().Format(time.RFC3339),
		TotalRegions: len(results),
		Results:      results,
	}

	var fastest, slowest *PingResult
	var fastestLatency, slowestLatency int64 = -1, -1

	for i := range results {
		if results[i].Success {
			summary.Successful++

			if fastestLatency == -1 || results[i].LatencyMS < fastestLatency {
				fastestLatency = results[i].LatencyMS
				fastest = &results[i]
			}

			if slowestLatency == -1 || results[i].LatencyMS > slowestLatency {
				slowestLatency = results[i].LatencyMS
				slowest = &results[i]
			}
		} else {
			summary.Failed++
		}
	}

	summary.Fastest = fastest
	summary.Slowest = slowest

	return summary
}
