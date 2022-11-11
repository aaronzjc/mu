package model

import "time"

const (
	AuthGithub = iota
	AuthWeibo
)

type User struct {
	ID       int       `gorm:"id"`
	Username string    `gorm:"username"`
	Nickname string    `gorm:"nickname"`
	Avatar   string    `gorm:"avatar"`
	AuthType string    `gorm:"auth_type"`
	AuthTime time.Time `gorm:"auth_time"`
	Token    string    `gorm:"token"`
	ExpireAt int64     `gorm:"expire_at"`
}

func (u *User) TableName() string {
	return "user"
}
