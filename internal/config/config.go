package config

import (
	"context"

	"github.com/go-redis/redis/v8"
)

var (
	rdb *redis.Client
	Ctx context.Context = context.Background()
)

func Init() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

}

func GetRedisDB() *redis.Client {
	return rdb
}
