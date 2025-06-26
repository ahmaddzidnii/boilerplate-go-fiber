//go:build wireinject
// +build wireinject

package injector

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/database"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/handlers"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/middlewares"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/routes"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type Application struct {
	App    *fiber.App
	Logger *logrus.Logger
}

func ProvideDatabase(logger *logrus.Logger) *gorm.DB {
	db := database.InitDatabase()
	logger.Info("Koneksi database dan migrasi berhasil.")
	return db
}

func ProvideRedis(logger *logrus.Logger) *redis.Client {
	client, err := database.InitRedis()
	if err != nil {
		logger.Error(err)
		logger.Panic("Gagal menginisialisasi Redis client")
	}
	logger.Info("Koneksi Redis berhasil.")
	return client
}

func ProvideValidator() *validator.Validate {
	return config.InitValidator()
}

func ProvideLogger() *logrus.Logger {
	return config.InitLogger()
}

func ProvideRouter(authHandler *handlers.AuthHandler, middleware *middlewares.Middleware, DB *gorm.DB) *fiber.App {
	app := fiber.New()
	routes.RegisterRoutes(app, authHandler, middleware, DB)
	return app
}

func NewApplication(app *fiber.App, logger *logrus.Logger) Application {
	return Application{
		App:    app,
		Logger: logger,
	}
}

var RepositorySet = wire.NewSet(
	//Mengembaikan constructor untuk setiap repository.
	repository.NewAuthRepository,
	repository.NewSessionRepository,
)

var ServiceSet = wire.NewSet(
	// Mengembalikan constructor untuk setiap service.
	service.NewAuthService,
)

var HandlerSet = wire.NewSet(
	handlers.NewAuthHandler,
	middlewares.NewMiddleware,
)

var AppSet = wire.NewSet(
	ProvideDatabase,
	ProvideRedis,
	ProvideValidator,
	ProvideLogger,
	ProvideRouter,
	NewApplication,

	// Set untuk setiap layer
	RepositorySet,
	ServiceSet,
	HandlerSet,
)

func InitializeApp() (Application, error) {
	wire.Build(AppSet)
	return Application{}, nil // Nilai dummy, akan diisi oleh Wire
}
