package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
)

type Middleware struct {
	Redis *redis.Client
}

func NewMiddleware(redis *redis.Client) *Middleware {
	return &Middleware{Redis: redis}
}

func (m *Middleware) AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var token string

		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		if token == "" {
			token = c.Cookies("session_id")
		}

		if token == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "Header otentikasi atau cookie sesi tidak ditemukan")
		}

		sessionKey := "session:" + token

		dataFromRedis, err := m.Redis.Get(c.Context(), sessionKey).Result()

		if err != nil {
			if errors.Is(err, redis.Nil) {
				utils.ClearCookies(c, "session_id")
				return utils.Error(c, fiber.StatusUnauthorized, "Sesi tidak ditemukan atau sudah berakhir")
			}
			log.Printf("Gagal mengambil sesi dari Redis: %v", err)
			return utils.Error(c, fiber.StatusUnauthorized, "Sesi tidak valid")
		}

		var sessionData models.Session
		if err := json.Unmarshal([]byte(dataFromRedis), &sessionData); err != nil {
			log.Printf("Data sesi corrupt di Redis: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
		}

		c.Locals("session_data", sessionData)
		return c.Next()
	}
}
