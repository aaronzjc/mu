package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"mu/internal/util/config"
	"mu/internal/util/logger"
	"sync"
)

var (
	client *redis.Client
	once sync.Once
)

func RedisConn() *redis.Client {
	if client == nil {
		once.Do(func() {
			cnf := config.NewConfig()
			client = redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d",cnf.Redis.Host, cnf.Redis.Port),
				Password: cnf.Redis.Password,
				DB:       0,
			})
		})
	}

	return client
}

func SaveToRedis(key string, hkey string, data string) {
	client := RedisConn()

	_, err := client.HSet(key, hkey, data).Result()
	if err != nil {
		logger.Error("SaveToRedis error , err = %s .", err.Error())
	}
}
