package db

import (
	"github.com/go-redis/redis"
)

var client *redis.Client

func InitRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})
}

func GetRedis() *redis.Client {
	return client
}
