// Package middleware provides JWT authentication middleware for Fiber.
package middleware

import (
	"strings"

	"challenge-api-go/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware creates a Fiber middleware that validates JWT tokens.
func JWTMiddleware(secret string) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "UNAUTHORIZED",
					Message: "Falta el header de autorización",
				},
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "UNAUTHORIZED",
					Message: "Formato de autorización inválido. Use: Bearer <token>",
				},
			})
		}

		tokenStr := parts[1]
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusForbidden).JSON(models.APIResponse{
				Success: false,
				Error: &models.APIError{
					Code:    "FORBIDDEN",
					Message: "Token inválido o expirado",
				},
			})
		}

		// Store the raw token for forwarding to downstream services.
		c.Locals("token", tokenStr)
		return c.Next()
	}
}
