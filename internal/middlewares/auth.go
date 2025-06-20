package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9" // <-- Import redis client
	"log"
)

// Struct Middleware akan menyimpan dependensi yang diperlukan
type Middleware struct {
	Redis *redis.Client
}

// Constructor untuk Middleware
func NewMiddleware(redis *redis.Client) *Middleware {
	return &Middleware{Redis: redis}
}

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("session_id")
		if sessionId == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "Token tidak diberikan")
		}

		sessionKey := "session:" + sessionId

		dataFromRedis, err := m.Redis.Get(c.Context(), sessionKey).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return utils.Error(c, fiber.StatusUnauthorized, "Session tidak ditemukan atau sudah berakhir")
			}
			log.Printf("Gagal mengambil session dari Redis: %v", err)
			return utils.Error(c, fiber.StatusUnauthorized, "Session tidak valid")
		}

		var sessionData models.Session
		if err := json.Unmarshal([]byte(dataFromRedis), &sessionData); err != nil {
			log.Printf("Data session corrupt injector Redis: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
		}

		c.Locals("session_data", sessionData)
		return c.Next()
	}
}
