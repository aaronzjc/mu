package model

import (
	"github.com/jinzhu/gorm"
	"mu/internal/util/config"
	"mu/internal/util/db"
	"mu/internal/util/logger"
)

var (
	Pool db.DB
)

type Query struct {
	Query string
	Args  []interface{}
	Order string
	Limit int
}

func DPool() db.DB {
	if Pool != (db.DB{}) {
		return Pool
	}

	c := config.NewConfig()
	Pool := db.DB{}
	Pool.Connect(&c)

	return Pool
}

func Prepare(q Query) *gorm.DB {
	conn := DPool().Conn

	if q.Order != "" {
		conn = conn.Order(q.Order)
	}

	if q.Query != "" {
		conn = conn.Where(q.Query, q.Args...)
	}

	if q.Limit > 0 {
		conn = conn.Limit(q.Limit)
	}

	return conn
}

func Select(fields string, q Query, model interface{}) error {
	conn := Prepare(q)
	defer conn.Close()

	conn.Select(fields).Find(model)
	if err := conn.Error; err != nil && !conn.RecordNotFound() {
		logger.Error("Select err %v, exp %s .", err, conn.QueryExpr())
		return err
	}

	return nil
}

func First(q Query, model interface{}) error {
	conn := Prepare(q)
	defer conn.Close()

	conn = conn.First(model)
	if err := conn.Error; err != nil && !conn.RecordNotFound() {
		logger.Error("First err %v, exp %s .", err, conn.QueryExpr())
		return err
	}

	return nil
}

func Update(model interface{}, data map[string]interface{}) error {
	conn := DPool().Conn
	defer conn.Close()

	conn = conn.Model(model).Update(data)
	if err := conn.Error; err != nil {
		logger.Error("Update err %v, exp %s .", err, conn.QueryExpr())
		return err
	}

	return nil
}

func Create(model interface{}) error {
	conn := DPool().Conn
	defer conn.Close()

	conn = conn.Create(model)
	if err := conn.Error; err != nil {
		logger.Error("Create err %v, exp %s .", err, conn.QueryExpr())
		return err
	}

	return nil
}

func Del(model interface{}) error {
	conn := DPool().Conn
	defer conn.Close()

	conn = conn.Delete(model)
	if err := conn.Error; err != nil {
		logger.Error("Del err %v, exp %s .", err, conn.QueryExpr())
		return err
	}

	return nil
}

func FetchRows(q Query, res interface{}) error {
	conn := Prepare(q)
	defer conn.Close()

	conn = conn.Find(res)
	if err := conn.Error; err != nil {
		logger.Error("FetchRows err %v, exp %s .", err, conn.QueryExpr())
		return err
	}

	return nil
}
