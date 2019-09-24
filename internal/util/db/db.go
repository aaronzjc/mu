package db

import (
	"crawler/internal/util/config"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var (
	Pool DB
)

type DB struct {
	Conn *gorm.DB
}

func (db *DB) Connect(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.Db.User,
		cfg.Db.Password,
		cfg.Db.Host,
		cfg.Db.Port,
		cfg.Db.Database,
	)
	c, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic("error connect to DB " + err.Error())
	}

	if gin.Mode() != gin.ReleaseMode {
		c.LogMode(true)
	}

	db.Conn = c
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	err := db.Conn.Close()

	if err != nil {
		log.Printf("[error] db close error %v \n", err)
		return err
	}

	return nil
}