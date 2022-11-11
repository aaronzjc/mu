package model

import (
	"time"
)

type Favor struct {
	ID        int       `gorm:"id"`
	UserId    int       `gorm:"user_id"`
	Site      string    `gorm:"site"`
	Key       string    `gorm:"key"`
	OriginUrl string    `gorm:"origin_url"`
	Title     string    `gorm:"title"`
	CreateAt  time.Time `gorm:"create_at"`
}

func (f *Favor) TableName() string {
	return "favor"
}
