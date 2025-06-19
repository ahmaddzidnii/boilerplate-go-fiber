package database

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func OpenConnetion() {
	host := config.GetEnv("DB_HOST", "localhost")
	user := config.GetEnv("DB_USER", "gorm")
	password := config.GetEnv("DB_PASSWORD", "gorm")
	dbname := config.GetEnv("DB_NAME", "gorm")
	port := config.GetEnv("DB_PORT", "9920")
	sslmode := config.GetEnv("DB_SSLMODE", "disable")

	dsn := "user=" + user + " password=" + password + " dbname=" + dbname + " host=" + host + " port=" + port + " sslmode=" + sslmode

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	}), &gorm.Config{})
	if err != nil {
		panic("Gagal membuka koneksi ke database: " + err.Error())
	}

	DB = db
}
