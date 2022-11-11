package model

import "time"

const DB_MU = "mu"

type BaseModel struct {
	ID        uint      `gorm:"column:id" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
