package utils

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

type Response struct {
	Status int         `json:"status"`
	Errors interface{} `json:"errors,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func Success(c *fiber.Ctx, statusCode int, data interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status: statusCode,
		Errors: nil,
		Data:   data,
	})
}

func Error(c *fiber.Ctx, statusCode int, errors interface{}) error {
	return c.Status(statusCode).JSON(Response{
		Status: statusCode,
		Errors: errors,
		Data:   nil,
	})
}

func ClearCookies(c *fiber.Ctx, key ...string) {
	for i := range key {
		c.Cookie(&fiber.Cookie{
			Name:    key[i],
			Expires: time.Now().Add(-time.Hour * 24),
			Value:   "",
		})
	}
}

func GetLocals[T any](c *fiber.Ctx, key string) (T, error) {
	var result T // Variabel penampung dengan tipe generik T

	data := c.Locals(key)
	if data == nil {
		return result, fmt.Errorf("data dengan key '%s' tidak ditemukan di context", key)
	}

	result, ok := data.(T)
	if !ok {
		// %T akan mencetak tipe data, berguna untuk debugging
		return result, fmt.Errorf("gagal konversi data: tipe data di context (%T) tidak cocok dengan tipe yang diminta (%T)", data, result)
	}

	return result, nil
}

func GenerateValidationResponse(err error) map[string]string {
	errorBag := make(map[string]string)
	// Type assertion ke validator.ValidationErrors
	for _, e := range err.(validator.ValidationErrors) {
		// e.Field()  -> Nama field struct (e.g., "Username")
		// e.Tag()    -> Aturan yang gagal (e.g., "required", "min")
		// e.Param()  -> Parameter aturan (e.g., "8" untuk "min=8")

		fieldName := strings.ToLower(e.Field())

		switch e.Tag() {
		case "required":
			errorBag[fieldName] = fmt.Sprintf("Field %s tidak boleh kosong.", fieldName)
		case "email":
			errorBag[fieldName] = fmt.Sprintf("Field %s harus berupa format email yang valid.", fieldName)
		case "min":
			errorBag[fieldName] = fmt.Sprintf("Field %s harus memiliki panjang minimal %s karakter.", fieldName, e.Param())
		case "max":
			errorBag[fieldName] = fmt.Sprintf("Field %s harus memiliki panjang maksimal %s karakter.", fieldName, e.Param())
		default:
			errorBag[fieldName] = fmt.Sprintf("Field %s tidak valid.", fieldName)
		}
	}
	return errorBag
}
