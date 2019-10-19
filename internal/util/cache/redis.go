package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	config2 "mu/internal/util/config"
	"mu/internal/util/logger"
)

func RedisConn() *redis.Client {
	c := config2.NewConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: "",
		DB:       0,
	})

	return client
}

func SaveToRedis(key string, hkey string, data string) {
	client := RedisConn()
	defer client.Close()

	_, err := client.HSet(key, hkey, data).Result()
	if err != nil {
		logger.Error("SaveToRedis error , err = %s .", err.Error())
	}
}
