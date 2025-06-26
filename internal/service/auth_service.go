package service

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/repository"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Username string `validate:"required" json:"username"` // Sesuai dengan request asli
	Password string `validate:"required" json:"password"`
}

// AuthService adalah interface untuk semua logika bisnis otentikasi
type AuthService interface {
	Login(ctx context.Context, req LoginRequest) (sessionID string, err error)
	Logout(ctx context.Context, sessionID string) error
}

type AuthServiceImpl struct {
	AuthRepository    repository.AuthRepository
	SessionRepository repository.SessionRepository
	Logger            *logrus.Logger
}

// NewAuthService adalah constructor untuk service
func NewAuthService(
	authRepo repository.AuthRepository,
	sessionRepo repository.SessionRepository,
	logger *logrus.Logger,
) AuthService {
	return &AuthServiceImpl{
		AuthRepository:    authRepo,
		SessionRepository: sessionRepo,
		Logger:            logger,
	}
}

var TTL = 2 * time.Hour // Durasi sesi, bisa diubah sesuai kebutuhan
var (
	ErrInvalidCredentials = errors.New("kombinasi NIM dan password salah")
	ErrInternalServer     = errors.New("terjadi kesalahan internal pada server")
)

func (s *AuthServiceImpl) Login(ctx context.Context, req LoginRequest) (string, error) {

	log := s.Logger.WithField("nim", req.Username)
	log.Info("Memproses permintaan login")

	// 2. Cari user via repository
	mhs, err := s.AuthRepository.FindByNIM(req.Username)
	if err != nil {
		// Jika user tidak ditemukan, kembalikan error umum
		log.Warn("Percobaan login gagal: NIM tidak ditemukan")
		return "", ErrInvalidCredentials
	}

	// 3. SEKARANG: Bandingkan password dengan aman menggunakan bcrypt
	// INI PERUBAHAN PALING PENTING!
	// Kode ini mengasumsikan password di DB Anda sudah di-hash dengan bcrypt.
	// Jika belum, Anda harus membuat mekanisme untuk hashing password saat registrasi.
	// err = bcrypt.CompareHashAndPassword([]byte(mhs.Password), []byte(req.Password))
	// if err != nil {
	// 	// Password tidak cocok
	// 	return "", errors.New("NIM atau password salah"), nil
	// }

	// SEMENTARA: Pakai perbandingan string biasa (TIDAK AMAN, ganti dengan bcrypt di atas)
	if mhs.Password != req.Password {
		log.Warn("Percobaan login gagal: Password salah")
		return "", ErrInvalidCredentials
	}

	// 4. Buat sesi
	sessionPayload := &models.Session{
		UserId: mhs.IdMahasiswa.String(),
		Nim:    mhs.NIM,
		Nama:   mhs.Nama,
	}
	sessionID := uuid.NewString()

	// 5. Simpan sesi ke Redis via repository
	err = s.SessionRepository.Create(ctx, sessionID, sessionPayload, TTL)
	if err != nil {
		log.WithError(err).Error("Gagal membuat sesi setelah otentikasi berhasil")
		return "", ErrInternalServer
	}

	// 6. Login berhasil, kembalikan sessionID
	log.WithField("session_id", sessionID).Info("Login berhasil, sesi dibuat")
	return sessionID, nil
}

func (s *AuthServiceImpl) Logout(ctx context.Context, sessionID string) error {
	s.Logger.WithField("session_id", sessionID).Info("Memproses permintaan logout")
	err := s.SessionRepository.Delete(ctx, sessionID)
	if err != nil {
		s.Logger.WithError(err).Error("Gagal menghapus sesi saat logout")
		return ErrInternalServer
	}
	return nil
}
