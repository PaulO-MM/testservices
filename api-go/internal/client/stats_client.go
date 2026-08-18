// Package client provides an HTTP client for calling the Node stats service.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"challenge-api-go/internal/models"
)

// StatsClient calls the Node.js stats API.
type StatsClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewStatsClient creates a client with a 5-second timeout.
func NewStatsClient(baseURL string) *StatsClient {
	return &StatsClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// GetStats sends Q and R matrices to the Node service and returns stats.
func (c *StatsClient) GetStats(q, r [][]float64, token string) (*models.MatrixStats, error) {
	reqBody := models.StatsRequest{}
	reqBody.Matrices.Q = q
	reqBody.Matrices.R = r

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stats request: %w", err)
	}

	url := c.BaseURL + "/api/v1/matrix/stats"
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("node service unavailable: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("node service returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var apiResp models.APIResponse
	if err := json.Unmarshal(respBody, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if !apiResp.Success {
		return nil, fmt.Errorf("node service returned success=false")
	}

	dataBytes, err := json.Marshal(apiResp.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	var stats models.MatrixStats
	if err := json.Unmarshal(dataBytes, &stats); err != nil {
		return nil, fmt.Errorf("failed to parse stats: %w", err)
	}

	return &stats, nil
}
