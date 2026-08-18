// Package handlers implements the QR decomposition HTTP endpoint.
package handlers

import (
	"fmt"
	"math"

	"challenge-api-go/internal/client"
	"challenge-api-go/internal/models"
	"challenge-api-go/internal/services"

	"github.com/gofiber/fiber/v2"
)

// Handler holds dependencies for the matrix endpoints.
type Handler struct {
	StatsClient *client.StatsClient
}

// NewHandler creates a new matrix handler.
func NewHandler(sc *client.StatsClient) *Handler {
	return &Handler{StatsClient: sc}
}

// ComputeQR handles POST /api/v1/matrix/qr.
func (h *Handler) ComputeQR(c *fiber.Ctx) error {
	var req models.QRRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_BODY",
				Message: "El cuerpo de la solicitud no es válido",
			},
		})
	}

	if len(req.Matrix) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_MATRIX",
				Message: "La matriz no puede estar vacía",
			},
		})
	}

	// Validate rectangular: all rows must have the same length.
	rowLengths := make([]int, len(req.Matrix))
	for i, row := range req.Matrix {
		rowLengths[i] = len(row)
		if len(row) == 0 {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_MATRIX",
					Message: "La matriz contiene filas vacías",
				},
			})
		}
	}
	for i := 1; i < len(rowLengths); i++ {
		if rowLengths[i] != rowLengths[0] {
			return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "INVALID_MATRIX",
					Message: "Las filas de la matriz tienen longitudes distintas",
					Details: map[string]interface{}{
						"rowLengths": rowLengths,
					},
				},
			})
		}
	}

	// Validate numeric elements.
	for i, row := range req.Matrix {
		for j, val := range row {
			if math.IsNaN(val) || math.IsInf(val, 0) {
				return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
					Success: false,
					Error: &models.APIError{
						Code:    "INVALID_MATRIX",
						Message: fmt.Sprintf("Elemento no numérico en posición [%d][%d]", i, j),
					},
				})
			}
		}
	}

	m := len(req.Matrix)
	n := len(req.Matrix[0])

	if m < n {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "UNSUPPORTED_DIMENSIONS",
				Message: "No se puede calcular QR: la matriz tiene más columnas que filas",
			},
		})
	}

	q, r, err := services.ComputeQR(req.Matrix)
	if err != nil {
		if err == services.ErrSingularMatrix {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "SINGULAR_MATRIX",
					Message: "La matriz tiene columnas linealmente dependientes",
				},
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Error al calcular la descomposición QR",
			},
		})
	}

	qrResult := &models.QRResult{
		Q: q,
		R: r,
		Dimensions: models.Dimensions{
			Rows: m,
			Cols: n,
		},
	}

	// Attempt to call Node stats service; degrade gracefully on failure.
	token, _ := c.Locals("token").(string)
	stats, err := h.StatsClient.GetStats(q, r, token)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(models.APIResponse{
			Success: true,
			Data: models.QRData{
				QR:    qrResult,
				Stats: nil,
			},
			Warning: fmt.Sprintf("No se pudieron calcular las estadísticas: %v", err),
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Data: models.QRData{
			QR:    qrResult,
			Stats: stats,
		},
	})
}
