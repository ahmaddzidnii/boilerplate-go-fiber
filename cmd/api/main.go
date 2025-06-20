package main

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/injector"
	"github.com/sirupsen/logrus"
)

func main() {
	config.LoadConfig()
	app, err := injector.InitializeApp()
	if err != nil {
		logrus.Fatal("Gagal menginisialisasi aplikasi: %v", err)
	}

	logger := app.Logger

	port := config.GetEnv("APP_PORT", "8080")

	logger.Info("🚀 Server siap dan berjalan injector port ", port)

	err = app.App.Listen("127.0.0.1:" + port)
	if err != nil {
		logger.Fatal("Gagal menjalankan server: ", err)
	}
}
