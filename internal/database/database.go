package database

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"time"
)

func InitDatabase() *gorm.DB {
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

	// =================================================================
	// PENGATURAN CONNECTION POOL
	// =================================================================
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Gagal mendapatkan objek sql.DB dari GORM: %v", err)
	}

	// 1. SetMaxIdleConns: Jumlah koneksi yang "diam" atau tidak terpakai di dalam pool.
	// Nilai yang baik adalah sama dengan MaxOpenConns untuk aplikasi dengan beban tinggi
	// agar tidak ada proses pembuatan koneksi baru saat ada lonjakan traffic.
	sqlDB.SetMaxIdleConns(100)

	// 2. SetMaxOpenConns: Jumlah maksimum koneksi yang boleh dibuka ke database.
	// Ini adalah 'rem' utama Anda untuk melindungi database.
	// Nilai 100-150 adalah titik awal yang kuat untuk 30k users.
	sqlDB.SetMaxOpenConns(120)

	// 3. SetConnMaxLifetime: Waktu maksimum sebuah koneksi boleh digunakan kembali.
	// Penting untuk load balancing di environment cloud (AWS RDS, Google Cloud SQL)
	// dan untuk menangani restart database dengan mulus.
	sqlDB.SetConnMaxLifetime(time.Hour) // 1 jam

	// 4. SetConnMaxIdleTime: Waktu maksimum sebuah koneksi boleh idle sebelum ditutup.
	// Ini membantu membebaskan resource di database selama periode sepi.
	// Harus lebih pendek dari ConnMaxLifetime.
	sqlDB.SetConnMaxIdleTime(10 * time.Minute) // 10 menit

	return db
}
