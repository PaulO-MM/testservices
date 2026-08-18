// Integration tests for the QR matrix HTTP handler.
package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"challenge-api-go/internal/auth"
	"challenge-api-go/internal/client"
	"challenge-api-go/internal/handlers"
	"challenge-api-go/internal/middleware"
	"challenge-api-go/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSecret = "test_secret_key"

// mockNodeServer creates a test HTTP server that mimics the Node stats service.
func mockNodeServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/api/v1/matrix/stats", r.URL.Path)

		authHeader := r.Header.Get("Authorization")
		assert.True(t, len(authHeader) > 7, "Authorization header should contain Bearer token")

		resp := models.APIResponse{
			Success: true,
			Data: models.MatrixStats{
				Max:     2.0,
				Min:     0.0,
				Average: 0.83,
				Sum:     5.0,
				Diagonal: models.DiagonalFlags{
					Q: false,
					R: true,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
}

// createTestApp sets up a Fiber app with all routes for testing.
func createTestApp(statsURL string) *fiber.App {
	app := fiber.New()

	statsClient := client.NewStatsClient(statsURL)
	authHandler := auth.NewHandler(testSecret)
	matrixHandler := handlers.NewHandler(statsClient)

	app.Post("/api/v1/auth/login", authHandler.Login)
	api := app.Group("/api/v1", middleware.JWTMiddleware(testSecret))
	api.Post("/matrix/qr", matrixHandler.ComputeQR)

	return app
}

// generateTestToken creates a valid JWT for testing.
func generateTestToken(t *testing.T) string {
	t.Helper()
	claims := jwt.MapClaims{
		"sub": "candidate",
		"exp": 9999999999,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(testSecret))
	require.NoError(t, err)
	return tokenStr
}

func TestMatrixHandler_NoJWT(t *testing.T) {
	mockServer := mockNodeServer(t)
	defer mockServer.Close()

	app := createTestApp(mockServer.URL)

	body, _ := json.Marshal(map[string]interface{}{
		"matrix": [][]float64{{1, 2}, {0, 1}, {1, 0}},
	})

	req := httptest.NewRequest("POST", "/api/v1/matrix/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var apiResp models.APIResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	assert.False(t, apiResp.Success)
	assert.Equal(t, "UNAUTHORIZED", apiResp.Error.Code)
}

func TestMatrixHandler_InvalidMatrix(t *testing.T) {
	mockServer := mockNodeServer(t)
	defer mockServer.Close()

	app := createTestApp(mockServer.URL)
	token := generateTestToken(t)

	// Empty matrix.
	body, _ := json.Marshal(map[string]interface{}{
		"matrix": [][]float64{},
	})

	req := httptest.NewRequest("POST", "/api/v1/matrix/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var apiResp models.APIResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	assert.False(t, apiResp.Success)
	assert.Equal(t, "INVALID_MATRIX", apiResp.Error.Code)
}

func TestMatrixHandler_ValidRequest(t *testing.T) {
	mockServer := mockNodeServer(t)
	defer mockServer.Close()

	app := createTestApp(mockServer.URL)
	token := generateTestToken(t)

	body, _ := json.Marshal(map[string]interface{}{
		"matrix": [][]float64{{1, 2}, {0, 1}, {1, 0}},
	})

	req := httptest.NewRequest("POST", "/api/v1/matrix/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var apiResp models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	require.NoError(t, err)
	assert.True(t, apiResp.Success)
	assert.NotEmpty(t, apiResp.Data)
	assert.Empty(t, apiResp.Warning, "Should have no warning when Node responds")
}

func TestMatrixHandler_UnsupportedDimensions(t *testing.T) {
	mockServer := mockNodeServer(t)
	defer mockServer.Close()

	app := createTestApp(mockServer.URL)
	token := generateTestToken(t)

	// m < n.
	body, _ := json.Marshal(map[string]interface{}{
		"matrix": [][]float64{{1, 2, 3}, {4, 5, 6}},
	})

	req := httptest.NewRequest("POST", "/api/v1/matrix/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusUnprocessableEntity, resp.StatusCode)

	var apiResp models.APIResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	assert.False(t, apiResp.Success)
	assert.Equal(t, "UNSUPPORTED_DIMENSIONS", apiResp.Error.Code)
}

func TestMatrixHandler_NodeFailureDegrades(t *testing.T) {
	// Mock server that always returns 500.
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"success":false,"error":{"code":"INTERNAL_ERROR","message":"fail"}}`)
	}))
	defer mockServer.Close()

	app := createTestApp(mockServer.URL)
	token := generateTestToken(t)

	body, _ := json.Marshal(map[string]interface{}{
		"matrix": [][]float64{{1, 2}, {0, 1}, {1, 0}},
	})

	req := httptest.NewRequest("POST", "/api/v1/matrix/qr", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var apiResp models.APIResponse
	json.NewDecoder(resp.Body).Decode(&apiResp)
	assert.True(t, apiResp.Success)
	assert.NotEmpty(t, apiResp.Warning, "Should warn about Node failure")
}
