package database

import (
	"context"
	"fmt"
	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitRedis() (*redis.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr := fmt.Sprintf("%s:%s", config.GetEnv("REDIS_URL", "localhost"), config.GetEnv("REDIS_PORT", "6379"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.GetEnv("REDIS_PASSWORD", ""),
		DB:       0,
		//TLSConfig: &tls.Config{},
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
