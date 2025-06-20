package routes

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/handlers"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/middlewares"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Ganti nama fungsi ini agar perannya lebih jelas (opsional tapi disarankan)
// Fungsi ini sekarang hanya bertugas mendaftarkan rute, bukan sebagai provider Wire.
func RegisterRoutes(app *fiber.App, authHandler *handlers.AuthHandler, mid *middlewares.Middleware) {
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://localhost:4173,https://situs-frontend-anda.com",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH",
		AllowCredentials: true,
	}))

	api := app.Group("/api/v1")

	authRoute := api.Group("/auth")
	authRoute.Post("/login", authHandler.Login)
	authRoute.Post("/logout", authHandler.Logout)
	authRoute.Get("/session", mid.AuthMiddleware(), authHandler.GetSession)
}
