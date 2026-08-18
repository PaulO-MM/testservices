// Main entry point for the QR decomposition API server.
package main

import (
	"log"

	"challenge-api-go/internal/auth"
	"challenge-api-go/internal/client"
	"challenge-api-go/internal/config"
	"challenge-api-go/internal/handlers"
	"challenge-api-go/internal/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	cfg := config.Load()

	app := fiber.New(fiber.Config{
		AppName: "QR Decomposition API",
	})

	// CORS middleware to allow frontend requests from different origins.
	app.Use(cors.New())

	statsClient := client.NewStatsClient(cfg.NodeAPIURL)
	authHandler := auth.NewHandler(cfg.JWTSecret)
	matrixHandler := handlers.NewHandler(statsClient)

	// Public routes.
	app.Post("/api/v1/auth/login", authHandler.Login)

	// Protected routes.
	api := app.Group("/api/v1", middleware.JWTMiddleware(cfg.JWTSecret))
	api.Post("/matrix/qr", matrixHandler.ComputeQR)

	log.Fatal(app.Listen(":" + cfg.Port))
}
