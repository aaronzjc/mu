package cache

import (
	config2 "crawler/internal/util/config"
	"crawler/internal/util/logger"
	"fmt"
	"github.com/go-redis/redis"
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
