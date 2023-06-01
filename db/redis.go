/**
 * Copyright (c) [2023] [Jim-shop]
 * [AircraftWarBackend] is licensed under Mulan PubL v2.
 * You can use this software according to the terms and conditions of the Mulan PubL v2.
 * You may obtain a copy of Mulan PubL v2 at:
 *          http://license.coscl.org.cn/MulanPubL-2.0
 * THIS SOFTWARE IS PROVIDED ON AN "AS IS" BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT,
 * MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PubL v2 for more details.
 */

package db

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
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
