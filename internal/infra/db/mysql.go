package db

import (
	"fmt"

	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/pkg/logger"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DBPool struct {
	dbMap map[string]*gorm.DB
}

var pool *DBPool

func init() {
	pool = &DBPool{
		dbMap: make(map[string]*gorm.DB),
	}
}

func Setup(conf *config.Config, config *gorm.Config) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("setup db panic")
		}
	}()
	var err error
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
		if pool.dbMap[dbname], err = gorm.Open(mysql.Open(dsn), config); err != nil {
			logger.Error("connect db " + dbname + " err")
			return err
		}
	}
	return nil
}

func Get(dbname string) (*gorm.DB, bool) {
	db, ok := pool.dbMap[dbname]
	return db, ok
}
