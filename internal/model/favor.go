package model

import (
	"errors"
	"time"
)

type Favor struct {
	ID			int			`gorm:"id"`
	UserId 		int 		`gorm:"user_id"`
	Site	 	string 		`gorm:"site"`
	Key 	   	string 		`gorm:"key"`
	OriginUrl 	string 		`gorm:"origin_url"`
	Title 		string 		`gorm:"title"`
	CreateAt 	time.Time 	`gorm:"create_at"`
}

type FavorJson struct {
	ID			int			`json:"id"`
	OriginUrl 	string 		`json:"origin_url"`
	Title 		string 		`json:"title"`
	CreateAt 	string 		`json:"create_at"`
}

func (f *Favor) TableName() string {
	return "favor"
}

func (f *Favor) Exist() (bool, error) {
	query := Query{}

	if f.ID > 0 {
		query.Query = "`id` = ?"
		query.Args = []interface{}{f.ID}
	} else if f.Key != "" {
		query.Query = "`user_id` = ? AND `site` = ? AND `key` = ?"
		query.Args = []interface{}{f.UserId, f.Site, f.Key}
	}

	err := First(query, &f)
	if err != nil {
		return false, errors.New("fetch error")
	}

	return f.ID > 0, nil
}

func (f *Favor) Create() error {
	err := Create(&f)
	if err != nil {
		return errors.New("create failed")
	}

	return nil
}

func (f *Favor) Del() bool {
	err := Del(&f)
	if err != nil {
		return false
	}

	return true
}

func (f *Favor) FetchRows(query Query) ([]Favor, error) {
	var list []Favor

	err := FetchRows(query, &list)
	if err != nil {
		return nil, errors.New("fetchRows favor failed")
	}

	return list, nil
}

func (f *Favor) Config(query Query) []string {
	var list []Favor

	err := Select("DISTINCT(`site`)", query, &list)
	if err != nil {
		return []string{}
	}

	var result []string
	for _, val := range list {
		result = append(result, val.Site)
	}

	return result
}

func (f *Favor) FormatJson() FavorJson {
	json := FavorJson{
		ID: f.ID,
		OriginUrl: f.OriginUrl,
		Title: f.Title,
		CreateAt: f.CreateAt.Format("2006-01-02 15:04:05"),
	}

	return json
}