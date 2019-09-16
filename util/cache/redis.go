package cache

import (
	"crawler/util/config"
	"fmt"
	"github.com/go-redis/redis"
	"log"
)

func RedisConn() *redis.Client {
	c := config.NewConfig()
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port),
		Password: "",
		DB:       0,
	})

	return client
}

func SaveToRedis(key string, hkey string, data string) {
	client := RedisConn()
	_, err := client.HSet(key, hkey, data).Result()
	if err != nil {
		log.Printf("[error] SaveToRedis error , err = %s\n", err.Error())
	}
}