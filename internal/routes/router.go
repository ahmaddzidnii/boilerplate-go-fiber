package routes

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/handlers"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/middlewares"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// Ganti nama fungsi ini agar perannya lebih jelas (opsional tapi disarankan)
// Fungsi ini sekarang hanya bertugas mendaftarkan rute, bukan sebagai provider Wire.

func RegisterRoutes(app *fiber.App, authHandler *handlers.AuthHandler, mid *middlewares.Middleware, DB *gorm.DB) {
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
	authRoute.Post("/logout", mid.AuthMiddleware(), authHandler.Logout)
	authRoute.Get("/session", mid.AuthMiddleware(), authHandler.GetSession)

	type MhsResponse struct {
		IDMahasiswa      string  `json:"id_mahasiswa"`
		NIM              string  `json:"nim"`
		Nama             string  `json:"nama"`
		IPK              float64 `json:"ipk"`
		IPSLalu          float64 `json:"ips_lalu"`
		TahunAkademik    string  `json:"tahun_akademik"`
		SemesterBerjalan int     `json:"semester_berjalan"`
		StatusMahasiswa  string  `json:"status_mahasiswa"`
		StatusPembayaran string  `json:"status_pembayaran"`
		CreatedAt        int64   `json:"created_at"`
		UpdatedAt        int64   `json:"updated_at"`
	}
	api.Get("/mhs", func(c *fiber.Ctx) error {
		var mhs []models.Mahasiswa

		if err := DB.Find(&mhs).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Gagal mengambil data mahasiswa",
			})
		}

		mhsResponses := make([]MhsResponse, 0, len(mhs))

		for _, m := range mhs {
			response := MhsResponse{
				IDMahasiswa:      m.IdMahasiswa.String(),
				NIM:              m.NIM,
				Nama:             m.Nama,
				IPK:              m.IPK,
				IPSLalu:          m.IPSLalu,
				TahunAkademik:    m.TahunAkademik,
				SemesterBerjalan: m.SemesterBerjalan,
				StatusMahasiswa:  m.StatusMahasiswa,
				StatusPembayaran: m.StatusPembayaran,
				CreatedAt:        m.CreatedAt.Unix(),
				UpdatedAt:        m.UpdatedAt.Unix(),
			}
			mhsResponses = append(mhsResponses, response)
		}

		return c.JSON(mhsResponses)
	})
}
