package model

import (
	"crawler/internal/util/config"
	"crawler/internal/util/db"
)

var (
	Pool db.DB
)

func DPool() db.DB {
	if Pool != (db.DB{}) {
		return Pool
	}

	c := config.NewConfig()
	Pool := db.DB{}
	Pool.Connect(&c)

	return Pool
}