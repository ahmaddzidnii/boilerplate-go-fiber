package repository

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"time"

	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/models"
	"github.com/redis/go-redis/v9"
)

// SessionRepository adalah interface untuk operasi data sesi di Redis
type SessionRepository interface {
	Create(ctx context.Context, sessionID string, payload *models.Session, ttl time.Duration) error
	Delete(ctx context.Context, sessionID string) error
	Get(ctx context.Context, sessionID string) (*models.Session, error)
}

type SessionRepositoryImpl struct {
	Redis  *redis.Client
	Logger *logrus.Logger
}

// NewSessionRepository adalah constructor
func NewSessionRepository(redis *redis.Client, logger *logrus.Logger) SessionRepository {
	return &SessionRepositoryImpl{
		Redis:  redis,
		Logger: logger,
	}
}

func (r *SessionRepositoryImpl) Create(ctx context.Context, sessionID string, payload *models.Session, ttl time.Duration) error {
	log := r.Logger.WithFields(logrus.Fields{
		"sessionID": sessionID,
		"userID":    payload.UserId,
		"nim":       payload.Nim,
	})

	log.Info("Menyimpan sesi ke Redis")

	// Marshal payload ke JSON di dalam repository
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.WithError(err).Error("Gagal mengubah payload ke JSON")
		return err
	}

	sessionKey := "session:" + sessionID
	errorRedis := r.Redis.Set(ctx, sessionKey, payloadBytes, ttl).Err()

	if errorRedis != nil {
		log.WithError(errorRedis).Error("Gagal menyimpan sesi ke Redis")
	}
	return errorRedis
}

func (r *SessionRepositoryImpl) Delete(ctx context.Context, sessionID string) error {
	log := r.Logger.WithField("session_id", sessionID)
	log.Info("Menghapus sesi dari Redis")

	sessionKey := "session:" + sessionID
	err := r.Redis.Del(ctx, sessionKey).Err()
	if err != nil {
		log.WithError(err).Error("Gagal menghapus sesi dari Redis")
	}
	return err
}

func (r *SessionRepositoryImpl) Get(ctx context.Context, sessionID string) (*models.Session, error) {
	log := r.Logger.WithField("session_id", sessionID)
	log.Info("Mengambil sesi dari Redis")

	sessionKey := "session:" + sessionID
	data, err := r.Redis.Get(ctx, sessionKey).Result()
	if err != nil {
		if err == redis.Nil {
			log.Warn("Sesi tidak ditemukan di Redis")
			return nil, nil // Sesi tidak ditemukan, bukan error
		}
		log.WithError(err).Error("Gagal mengambil sesi dari Redis")
		return nil, err // Error lain saat mengambil data
	}

	var session models.Session
	if err := json.Unmarshal([]byte(data), &session); err != nil {
		log.WithError(err).Error("Gagal mengubah data sesi dari JSON")
		return nil, err // Error saat unmarshal
	}

	return &session, nil
}
