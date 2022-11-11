package dto

import (
	"github.com/aaronzjc/mu/internal/domain/model"
	"github.com/aaronzjc/mu/pkg/helper"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	AuthType string `json:"auth_type"`
	AuthTime string `json:"auth_time"`
}

func (u *User) FillByModel(user model.User) *User {
	u.ID = user.ID
	u.Username = user.Username
	u.Nickname = user.Nickname
	u.Avatar = user.Avatar
	u.AuthType = user.AuthType
	u.AuthTime = helper.TimeToLocalStr(user.AuthTime)
	return u
}
