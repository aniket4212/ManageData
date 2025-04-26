package redis

import (
	"context"
	"log"
	"strings"

	"managedata/config"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func ConnectRedis() {
	redisConf := config.AppConfig.RedisConfig

	RDB = redis.NewClient(&redis.Options{
		Addr:     strings.TrimSpace(redisConf.Addr),
		Password: redisConf.Password,
		DB:       redisConf.DB,
	})

	if err := RDB.Ping(Ctx).Err(); err != nil {
		log.Fatalf("Error while connecting to Redis: %v", err)
	}

	log.Println("Redis Connected")
}
