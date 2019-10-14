package model

import (
	"crawler/internal/util/config"
	"crawler/internal/util/db"
	"crawler/internal/util/logger"
	"github.com/jinzhu/gorm"
)

var (
	Pool db.DB
)

type Query struct {
	Query string
	Args []interface{}
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
	db := DPool().Conn

	if q.Order != "" {
		db = db.Order(q.Order)
	}

	if q.Query != "" {
		db = db.Where(q.Query, q.Args...)
	}

	if q.Limit > 0 {
		db = db.Limit(q.Limit)
	}

	return db
}

func Select(fields string, q Query, model interface{}) error {
	db := Prepare(q)
	defer db.Close()

	db.Select(fields).Find(model)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		logger.Error("Select err %v, exp %s .", err, db.QueryExpr())
		return err
	}

	return nil
}

func First(q Query, model interface{}) error {
	db := Prepare(q)
	defer db.Close()

	db = db.First(model)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		logger.Error("First err %v, exp %s .", err, db.QueryExpr())
		return err
	}

	return nil
}

func Update(model interface{}, data map[string]interface{}) error {
	db := DPool().Conn
	defer db.Close()

	db = db.Model(model).Update(data)
	if err := db.Error; err != nil {
		logger.Error("Update err %v, exp %s .", err, db.QueryExpr())
		return err
	}

	return nil
}

func Create(model interface{}) error {
	db := DPool().Conn
	defer db.Close()

	db = db.Create(model)
	if err := db.Error; err != nil {
		logger.Error("Create err %v, exp %s .", err, db.QueryExpr())
		return err
	}

	return nil
}

func Del(model interface{}) error {
	db := DPool().Conn
	defer db.Close()

	db = db.Delete(model)
	if err := db.Error; err != nil {
		logger.Error("Del err %v, exp %s .", err, db.QueryExpr())
		return err
	}

	return nil
}

func FetchRows(q Query, res interface{}) error {
	db := Prepare(q)
	defer db.Close()

	db = db.Find(res)
	if err := db.Error; err != nil {
		logger.Error("FetchRows err %v, exp %s .", err, db.QueryExpr())
		return err
	}

	return nil
}