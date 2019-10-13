package model

import (
	"crawler/internal/util/cache"
	"crawler/internal/util/tool"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const Key = "admin_auth_%s"

type Admin struct {
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	AuthType int8   `json:"auth_type"`
	AuthTime string `json:"auth_time"`
	Token    string `json:"token"`
}

func (u *Admin) FetchInfo() error {
	conn := cache.RedisConn()
	defer conn.Close()

	res, err := conn.Get(u.AuthKey()).Result()
	if err != nil {
		return errors.New("fetch redis info failed")
	}

	_ = json.Unmarshal([]byte(res), u)

	return nil
}

func (u *Admin) Auth() {
	conn := cache.RedisConn()
	defer conn.Close()

	u.Token = tool.GenerateToken(u.Username)
	u.AuthTime = time.Now().Format("2006-01-02 15:04:05")

	info, _ := json.Marshal(u)
	conn.Set(u.AuthKey(), string(info), time.Hour * 24)
}

func (u *Admin) CheckToken(token string) bool {
	err := u.FetchInfo()
	if err != nil {
		return false
	}

	return u.Token == token
}

func (u *Admin) AuthKey() string {
	return fmt.Sprintf(Key, u.Username)
}