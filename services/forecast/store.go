package forecast

//This file will take the data from the met office api, extract the data it needs and store it within redis

import (
	"github.com/redis/go-redis/v9"
)

func redisClient() *redis.Client {

	rc := redis.NewClient(&redis.Options{
		Addr:     "redis-19441.c233.eu-west-1-1.ec2.redns.redis-cloud.com:19441",
		Password: "",
		DB:       0,
	})

	return rc

}
