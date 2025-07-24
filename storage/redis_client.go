package storage

import (
	"context"

	"github.com/redis/go-redis/v9"
)

var (
	Ctx = context.Background()
	RDB *redis.Client
)

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}
