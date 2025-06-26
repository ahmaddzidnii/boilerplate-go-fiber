package repository

import (
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"time"
)

// AuthRepository interface
type AuthRepository interface {
	FindByNIM(nim string) (*models.Mahasiswa, error)
}

type AuthRepositoryImpl struct {
	Db     *gorm.DB
	Logger *logrus.Logger
}

// NewAuthRepository constructor
func NewAuthRepository(db *gorm.DB, logger *logrus.Logger) AuthRepository {
	return &AuthRepositoryImpl{
		Db:     db,
		Logger: logger,
	}
}

func (r *AuthRepositoryImpl) FindByNIM(nim string) (*models.Mahasiswa, error) {
	log := r.Logger.WithField("nim", nim)
	log.Info("Mencari mahasiswa berdasarkan NIM di database")
	var mhs models.Mahasiswa
	startTime := time.Now()
	err := r.Db.Where("nim = ?", nim).First(&mhs).Error
	duration := time.Since(startTime)

	if duration > 5*time.Second {
		log.WithField("duration", duration).Warn("Pencarian mahasiswa dengan NIM memakan waktu lebih dari 5 detik")
	} else {
		log.WithField("duration", duration).Info("Pencarian mahasiswa selesai")
	}

	if err != nil {
		return nil, err
	}
	return &mhs, nil
}
