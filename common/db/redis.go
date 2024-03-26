package db

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"tiktok/common/config"
)

type Redis struct {
	Client *redis.Client
}

var rds = &Redis{}

func RedisInit() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Get().Redis.Address,
		Password: config.Get().Redis.Password,
		DB:       config.Get().Redis.Db,
	})

	if _, err := rdb.Ping(context.Background()).Result(); err != nil {
		log.Printf("connect redis failed: %v", err)
	}
	rds = &Redis{Client: rdb}
	log.Println("Connect redis succeeded")
}

func GetRedis() *redis.Client {
	return rds.Client
}
