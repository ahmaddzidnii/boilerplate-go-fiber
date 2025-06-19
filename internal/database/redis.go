package database

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"

	"github.com/ahmaddzidnii/backend-krs-auth-service/internal/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client
var RedisCtx = context.Background()

// OpenRedisConnection diubah untuk mengembalikan error jika koneksi gagal.
func OpenRedisConnection() {
	addr := fmt.Sprintf("%s:%s", config.GetEnv("REDIS_URL", "localhost"), config.GetEnv("REDIS_PORT", "6379"))

	rdb := redis.NewClient(&redis.Options{
		Addr:      addr,
		Password:  config.GetEnv("REDIS_PASSWORD", ""),
		DB:        0,
		TLSConfig: &tls.Config{},
	})

	if err := rdb.Ping(RedisCtx).Err(); err != nil {
		panic(err)

	}
	log.Println("Berhasil terhubung ke Redis.")
	RedisClient = rdb
}
