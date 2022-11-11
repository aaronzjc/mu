package cache

import (
	"errors"
	"fmt"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/go-redis/redis"
)

var (
	client *redis.Client
)

func Setup(conf *config.RedisConfig) error {
	if client == nil {
		client = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
			Password: conf.Password,
			DB:       0,
		})
		if _, err := client.Ping().Result(); err != nil {
			return errors.New("init redis err")
		}
	}
	return nil
}

func Get() *redis.Client {
	return client
}
