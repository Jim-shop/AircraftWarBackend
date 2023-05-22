package db

import (
	"github.com/spf13/viper"
	"github.com/go-redis/redis"
)

var rds *redis.Client

func InitRedis() {
	rds = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})
}

func GetRedis() *redis.Client {
	return rds
}

//todo
