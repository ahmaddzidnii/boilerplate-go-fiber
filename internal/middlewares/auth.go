package middlewares

import (
	"encoding/json"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/database"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"log"
)

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		sessionId := c.Cookies("session_id")
		if sessionId == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "Akses ditolak")
		}

		sessionKey := "session:" + sessionId
		dataFromRedis, err := database.RedisClient.Get(database.RedisCtx, sessionKey).Result()
		if err != nil {
			// ... (handle error redis.Nil dan error lainnya)
			return utils.Error(c, fiber.StatusUnauthorized, "Session tidak valid")
		}

		var sessionData models.Session
		if err := json.Unmarshal([]byte(dataFromRedis), &sessionData); err != nil {
			log.Printf("Data session corrupt di Redis: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "Server error")
		}

		// BERHASIL! Simpan STRUCT ke c.Locals
		c.Locals("session_data", sessionData)
		return c.Next()
	}
}
