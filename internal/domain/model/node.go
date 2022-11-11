package model

import "time"

const (
	PingFail int8 = 0
	PingOk   int8 = 1
)

type Node struct {
	ID       int       `gorm:"id" json:"id"`
	Name     string    `gorm:"name" json:"name"`
	Addr     string    `gorm:"addr" json:"addr"`
	Type     int8      `gorm:"type" json:"type"`
	Enable   int8      `gorm:"enable" json:"enable"`
	Ping     int8      `gorm:"ping" json:"ping"`
	CreateAt time.Time `gorm:"create_at" json:"create_at"`
}

func (node *Node) TableName() string {
	return "node"
}
