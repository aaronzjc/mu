package db

import (
	"errors"
	"fmt"
	"sync"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBPool struct {
	dbMap map[string]*gorm.DB
}

var (
	pool *DBPool
	once sync.Once
)

func init() {
	pool = &DBPool{
		dbMap: make(map[string]*gorm.DB),
	}
}

func Setup(conf *config.Config, config *gorm.Config) error {
	var err error
	once.Do(func() {
		for dbname, v := range conf.Database {
			// 初始化DB等
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True",
				v.Username,
				v.Password,
				v.Host,
				v.Port,
				dbname,
				v.Charset,
			)
			if client, errr := gorm.Open(mysql.Open(dsn), config); errr != nil {
				logger.Error("init db err, ", errr.Error())
				err = errors.New("connect to " + dbname + " err")
				return
			} else {
				pool.dbMap[dbname] = client
			}
		}
	})
	return err
}

func Get(dbname string) (*gorm.DB, bool) {
	db, ok := pool.dbMap[dbname]
	return db, ok
}
