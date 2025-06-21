package handlers

import (
	"encoding/json"
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"log"
	"time"
)

type AuthHandler struct {
	DB        *gorm.DB
	Redis     *redis.Client
	Validator *validator.Validate
}

func NewAuthHandler(db *gorm.DB, redis *redis.Client, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		DB:        db,
		Redis:     redis,
		Validator: validator,
	}
}

type LoginRequest struct {
	Username string `validate:"required" json:"username"`
	Password string `validate:"required" json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if err := h.Validator.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorBag := utils.GenerateValidationResponse(validationErrors)
			return utils.Error(c, fiber.StatusBadRequest, errorBag)
		}
	}

	mhs := models.Mahasiswa{}
	data := h.DB.First(&mhs, "nim = ? AND password = ?", req.Username, req.Password)

	if data.Error != nil {
		return utils.Error(c, fiber.StatusUnauthorized, "Username atau password salah")
	}

	sessionPayload := models.Session{
		UserId: mhs.IdMahasiswa,
		Nim:    mhs.Nim,
		Nama:   mhs.Nama,
	}
	sessionId := uuid.NewString()

	payload, errMarshal := json.Marshal(sessionPayload)
	if errMarshal != nil {
		log.Print(errMarshal)
		return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
	}

	sessionKey := "session:" + sessionId
	ttl := 2 * time.Hour

	errRedis := h.Redis.Set(c.Context(), sessionKey, payload, ttl).Err()
	if errRedis != nil {
		log.Printf("Gagal menyimpan session ke Redis: %v", errRedis)
		return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(ttl)
	cookie.HTTPOnly = true
	cookie.Path = "/"
	cookie.SameSite = "none"
	cookie.Secure = false

	c.Cookie(cookie)

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"session_id": sessionId,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	sessionId := c.Cookies("session_id")
	if sessionId == "" {
		return utils.Error(c, fiber.StatusUnauthorized, "Tidak ada session yang ditemukan")
	}

	sessionKey := "session:" + sessionId

	err := h.Redis.Del(c.Context(), sessionKey).Err()
	if err != nil {
		return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
	}

	utils.ClearCookies(c, "session_id")
	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"message": "Logout successful",
	})
}

func (h *AuthHandler) GetSession(c *fiber.Ctx) error {
	sessionData, err := utils.GetLocals[models.Session](c, "session_data")
	if err != nil {
		log.Printf("Gagal mendapatkan session data: %v", err)
		return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
	}

	return utils.Success(c, fiber.StatusOK, sessionData)
}
