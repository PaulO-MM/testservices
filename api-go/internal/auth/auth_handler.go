// Package auth handles the mock login endpoint for JWT token generation.
package auth

import (
	"time"

	"challenge-api-go/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// DECISION: Hardcoded mock credentials as specified. In production this would
// be an external identity service with proper password hashing and validation.
const (
	mockUsername = "candidate"
	mockPassword = "challenge2024"
	tokenExpiry  = 3600 // seconds
)

// Handler holds dependencies for the auth handlers.
type Handler struct {
	JWTSecret string
}

// NewHandler creates a new auth handler.
func NewHandler(jwtSecret string) *Handler {
	return &Handler{JWTSecret: jwtSecret}
}

// Login handles POST /api/v1/auth/login with mock credentials.
func (h *Handler) Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_BODY",
				Message: "El cuerpo de la solicitud no es válido",
			},
		})
	}

	if req.Username != mockUsername || req.Password != mockPassword {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Credenciales incorrectas",
			},
		})
	}

	claims := jwt.MapClaims{
		"sub": req.Username,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Duration(tokenExpiry) * time.Second).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Success: false,
			Error: &models.APIError{
				Code:    "INTERNAL_ERROR",
				Message: "Error al generar el token",
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Success: true,
		Data: models.LoginData{
			Token:     tokenStr,
			ExpiresIn: tokenExpiry,
		},
	})
}
