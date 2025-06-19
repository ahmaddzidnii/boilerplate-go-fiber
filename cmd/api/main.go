package main

import (
	"log"

	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/database"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()
	database.OpenConnetion()
	database.OpenRedisConnection()
	config.InitializeValidator()

	errorMigrate := database.DB.AutoMigrate(&models.Mahasiswa{})

	if errorMigrate != nil {
		log.Printf(errorMigrate.Error())
	}
	app := fiber.New()

	routes.SetupRoutes(app)

	port := config.GetEnv("APP_PORT", "8080")

	log.Printf("🚀 Server berjalan di port %s", port)

	err := app.Listen("127.0.0.1:" + port)

	if err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
