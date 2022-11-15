package test

import (
	"os"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/internal/infra/cache"
	"github.com/aaronzjc/mu/internal/infra/db"

	"gorm.io/gorm"
)

var (
	conf = config.Config{
		Database: map[string]config.DbConfig{
			"mu": {
				Host:     "127.0.0.1",
				Port:     3306,
				Username: "root",
				Password: "123456",
				Charset:  "utf8",
			},
		},
		Redis: config.RedisConfig{
			Host: "127.0.0.1",
			Port: 7379,
		},
	}
)

func SetupDb() error {
	return db.Setup(&conf, &gorm.Config{})
}

func SetupCache() error {
	return cache.Setup(&conf.Redis)
}

func SetupProxy() {
	os.Setenv("HTTP_PROXY", "http://127.0.0.1:51081")
	os.Setenv("HTTPS_PROXY", "http://127.0.0.1:51081")
}

func ClearProxy() {
	os.Unsetenv("HTTP_PROXY")
	os.Unsetenv("HTTPS_PROXY")
}
