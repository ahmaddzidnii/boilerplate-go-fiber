package handlers

import (
	"errors"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/service"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"log"
	"strings"
	"time"
)

type AuthHandler struct {
	AuthService service.AuthService
	Logger      *logrus.Logger
	Validator   *validator.Validate
}

func NewAuthHandler(authService service.AuthService, logger *logrus.Logger, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		AuthService: authService,
		Logger:      logger,
		Validator:   validator,
	}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req service.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		h.Logger.WithError(err).Error("Gagal mem-parsing body request login")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "fail", "message": err.Error()})
	}

	if err := h.Validator.Struct(req); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			errorBag := utils.GenerateValidationResponse(validationErrors)
			return utils.Error(c, fiber.StatusBadRequest, errorBag)
		}
	}

	sessionId, err := h.AuthService.Login(c.Context(), req)

	if err != nil {
		h.Logger.WithError(err).Error("Gagal melakukan login")
		if errors.Is(err, service.ErrInvalidCredentials) {
			return utils.Error(c, fiber.StatusUnauthorized, "NIM atau password salah")
		} else if errors.Is(err, service.ErrInternalServer) {
			return utils.Error(c, fiber.StatusInternalServerError, "Internal server error")
		}
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "session_id"
	cookie.Value = sessionId
	cookie.Expires = time.Now().Add(service.TTL)
	cookie.HTTPOnly = true
	cookie.Path = "/"

	c.Cookie(cookie)

	return utils.Success(c, fiber.StatusOK, fiber.Map{
		"session_id": sessionId,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	var sessionId string
	authHeader := c.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		sessionId = strings.TrimPrefix(authHeader, "Bearer ")

	} else {
		sessionId = c.Cookies("session_id")

	}

	if sessionId == "" {
		return utils.Error(c, fiber.StatusUnauthorized, "Tidak ada sesi atau token yang ditemukan")
	}

	err := h.AuthService.Logout(c.Context(), sessionId)

	if err != nil {
		h.Logger.WithError(err).Error("Gagal menghapus sesi")
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
