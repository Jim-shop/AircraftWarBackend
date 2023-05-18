package db

import (
	"github.com/go-redis/redis"
)

var rds *redis.Client

func InitRedis() {
	rds = redis.NewClient(&redis.Options{
		Addr:     "",
		Password: "",
		DB:       0,
	})
}

func GetRedis() *redis.Client {
	return rds
}
