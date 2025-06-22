package middlewares

import (
	"encoding/json"
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9" // <-- Import redis client
	"log"
	"strings"
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
		var token string

		// 1. Prioritaskan untuk memeriksa 'Authorization: Bearer <token>' header
		authHeader := c.Get("Authorization")
		if strings.HasPrefix(authHeader, "Bearer ") {
			token = strings.TrimPrefix(authHeader, "Bearer ")
		}

		// 2. Jika tidak ada di header, fallback ke cookie
		if token == "" {
			token = c.Cookies("session_id")
		}

		// 3. Jika token tetap kosong setelah dicek di kedua tempat, tolak akses
		if token == "" {
			return utils.Error(c, fiber.StatusUnauthorized, "Header otentikasi atau cookie sesi tidak ditemukan")
		}

		// --- Mulai dari sini, logikanya sama persis, hanya menggunakan variabel 'token' ---

		// 4. Verifikasi token/session ID dengan Redis
		sessionKey := "session:" + token
		dataFromRedis, err := m.Redis.Get(c.Context(), sessionKey).Result()
		if err != nil {
			// Jika session tidak ditemukan di Redis
			if errors.Is(err, redis.Nil) {
				utils.ClearCookies(c, "session_id")
				return utils.Error(c, fiber.StatusUnauthorized, "Sesi tidak ditemukan atau sudah berakhir")
			}
			// Jika ada error lain saat mengakses Redis
			log.Printf("Gagal mengambil sesi dari Redis: %v", err)
			return utils.Error(c, fiber.StatusUnauthorized, "Sesi tidak valid")
		}

		// 5. Unmarshal data sesi dan simpan ke c.Locals
		var sessionData models.Session
		if err := json.Unmarshal([]byte(dataFromRedis), &sessionData); err != nil {
			log.Printf("Data sesi corrupt di Redis: %v", err)
			return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
		}

		c.Locals("session_data", sessionData)

		// 6. Lanjutkan ke handler/route selanjutnya
		return c.Next()
	}
}
