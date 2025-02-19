package Rdb

import (
	"os"

	"github.com/go-redis/redis/v8"
)

func GetRedisDB() *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-17171.c338.eu-west-2-1.ec2.redns.redis-cloud.com:17171",
		Password: os.Getenv("redis"),
		DB:       0,
	})
	return rc
}
