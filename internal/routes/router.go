package routes

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/handlers"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	// Middleware untuk logging setiap request
	app.Use(logger.New())

	// ========================================================================
	// KONFIGURASI CORS
	// ========================================================================
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:4173,https://situs-frontend-anda.com",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH",
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")

	authRoute := api.Group("/auth")

	authRoute.Post("/login", handlers.Login)
	authRoute.Post("/logout", handlers.Logout)
	authRoute.Get("/session", middlewares.AuthMiddleware(), handlers.GetSession)

}
