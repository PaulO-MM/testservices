// Package models defines request/response types for the QR matrix API.
package models

// Standard API response envelope.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Warning string      `json:"warning,omitempty"`
}

// APIError represents a structured error returned to the client.
type APIError struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// QRRequest is the incoming payload for the QR decomposition endpoint.
type QRRequest struct {
	Matrix [][]float64 `json:"matrix"`
}

// Dimensions stores the row/column counts of a matrix.
type Dimensions struct {
	Rows int `json:"rows"`
	Cols int `json:"cols"`
}

// QRResult holds the Q and R matrices plus their dimensions.
type QRResult struct {
	Q          [][]float64 `json:"q"`
	R          [][]float64 `json:"r"`
	Dimensions Dimensions  `json:"dimensions"`
}

// DiagonalFlags indicates whether each matrix is diagonal.
type DiagonalFlags struct {
	Q bool `json:"q"`
	R bool `json:"r"`
}

// MatrixStats holds the statistical summary computed by the Node service.
type MatrixStats struct {
	Max      float64       `json:"max"`
	Min      float64       `json:"min"`
	Average  float64       `json:"average"`
	Sum      float64       `json:"sum"`
	Diagonal DiagonalFlags `json:"diagonal"`
}

// QRData combines QR decomposition output with optional stats.
type QRData struct {
	QR    *QRResult     `json:"qr"`
	Stats *MatrixStats  `json:"stats"`
}

// StatsRequest is the payload sent to the Node stats service.
type StatsRequest struct {
	Matrices struct {
		Q [][]float64 `json:"q"`
		R [][]float64 `json:"r"`
	} `json:"matrices"`
}

// LoginRequest represents the credentials payload for the mock login endpoint.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginData holds the JWT token and its expiry.
type LoginData struct {
	Token     string `json:"token"`
	ExpiresIn int    `json:"expiresIn"`
}
