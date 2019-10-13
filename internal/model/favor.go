package model

import (
	"crawler/internal/util/logger"
	"errors"
	"time"
)

type Favor struct {
	ID			int			`gorm:"id" json:"id"`
	UserId 		int 		`gorm:"user_id" json:"user_id"`
	Site	 	string 		`gorm:"site" json:"site"`
	Key 	   	string 		`gorm:"key" json:"key"`
	OriginUrl 	string 		`gorm:"origin_url" json:"origin_url"`
	Title 		string 		`gorm:"title" json:"title"`
	CreateAt 	time.Time 	`gorm:"create_at" json:"create_at"`
}

func (f *Favor) TableName() string {
	return "favor"
}

func (f *Favor) Exist() bool {
	db := DPool().Conn
	defer db.Close()

	if f.ID > 0 {
		db = db.Where("`id` = ?", f.ID).First(f)
	} else if f.Key != "" {
		db = db.Where(" `user_id` = ? AND `site` = ? AND `key` = ?", f.UserId, f.Site, f.Key).First(f)
	}
	if err := db.Error; err != nil && !db.RecordNotFound() {
		logger.Error("Exist err %v, exp %s .", err, db.QueryExpr())
		return true
	}
	return f.ID > 0
}

func (f *Favor) Create() error {
	db := DPool().Conn
	defer db.Close()

	db = db.Create(&f)
	var err error
	if err = db.Error; err != nil {
		logger.Error("create err %v, exp %s .", err, db.QueryExpr())
		return errors.New("create favor err")
	}

	return nil
}

func (f *Favor) Del() bool {
	db := DPool().Conn
	defer db.Close()

	db = db.Delete(f)
	if err := db.Error; err != nil {
		logger.Error("delete err %v, exp %s .", err, db.QueryExpr())
		return false
	}

	return true
}

func (f *Favor) FetchRows(query string, args ...interface{}) ([]Favor, error) {
	db := DPool().Conn
	defer db.Close()

	var list []Favor
	db = db.Where(query, args...).Find(&list)
	if err := db.Error; err != nil {
		logger.Error("FetchRows err %v, exp %s .", err, db.QueryExpr())
		return nil, errors.New("fetchRows favor failed")
	}
	return list, nil
}

func (f *Favor) Config() []string {
	db := DPool().Conn
	defer db.Close()

	var list []Favor
	db = db.Select("DISTINCT(`site`)").Find(&list)
	if err := db.Error; err != nil {
		logger.Error("FetchConfig err %v, exp %s .", err, db.QueryExpr())
		return []string{}
	}

	var result []string
	for _, val := range list {
		result = append(result, val.Site)
	}

	return result
}