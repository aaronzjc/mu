package model

import (
	"errors"
	"fmt"
	"mu/internal/util/logger"
	"mu/internal/util/tool"
	"time"
)

const (
	AuthGithub = iota
	AuthWeibo
)

type User struct {
	ID       int       `gorm:"id"`
	Username string    `gorm:"username"`
	Nickname string    `gorm:"nickname"`
	Avatar   string    `gorm:"avatar"`
	AuthType int8      `gorm:"auth_type"`
	AuthTime time.Time `gorm:"auth_time"`
	Token    string    `gorm:"token"`
	ExpireAt int64     `gorm:"expire_at"`
}

type UserJson struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	AuthType int8   `json:"auth_type"`
	AuthTime string `json:"auth_time"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) FetchRow(query Query) (User, error) {
	var tmp User
	err := First(query, &tmp)
	if err != nil {
		return User{}, errors.New("fetch user info failed")
	}

	return tmp, nil
}

func (u *User) Create() error {
	tmp, err := u.FetchRow(Query{
		Query: "`username` = ?",
		Args:  []interface{}{u.Username},
	})
	if err != nil {
		return errors.New("create user err")
	}

	if tmp.ID > 0 {
		return errors.New(fmt.Sprintf("user with %s exists", u.Username))
	}

	err = Create(&u)
	if err != nil {
		return errors.New("create user err")
	}

	return nil
}

func (u *User) Update(data map[string]interface{}) error {
	db := DPool().Conn
	defer db.Close()

	db = db.Model(&u).Update(data)
	if err := db.Error; err != nil {
		logger.Error("update err %v, exp %s .", err, db.QueryExpr())
		return errors.New("update user failed")
	}

	return nil
}

func (u *User) RefreshToken() error {
	token := tool.GenerateToken(u.Username)
	expireAt := time.Now().Add(time.Hour * 24 * 30).Unix()

	if err := u.Update(map[string]interface{}{
		"token":     token,
		"expire_at": expireAt,
	}); err != nil {
		return errors.New("refresh token error " + err.Error())
	}

	u.Token = token
	u.ExpireAt = expireAt

	return nil
}

func (u *User) Auth() error {
	tmp, _ := u.FetchRow(Query{
		Query: "`username` = ? AND `auth_type` = ?",
		Args:  []interface{}{u.Username, u.AuthType},
	})
	if tmp.ID > 0 {
		err := tmp.RefreshToken()
		if err != nil {
			return errors.New("auth failed " + err.Error())
		}
		u.ID = tmp.ID
		u.AuthTime = time.Now()
		u.Token = tmp.Token
		u.ExpireAt = tmp.ExpireAt

		err = u.Update(map[string]interface{}{
			"auth_time": u.AuthTime,
		})
		if err != nil {
			logger.Error("user auth err e = ", u.ID, err.Error())
		}
	} else {
		u.AuthTime = time.Now()
		u.Token = tool.GenerateToken(u.Username)
		u.ExpireAt = time.Now().Add(time.Hour * 24 * 30).Unix()

		err := u.Create()
		if err != nil {
			return errors.New("auth failed " + err.Error())
		}
	}

	return nil
}

func (u *User) CheckToken() (bool, error) {
	login, _ := u.FetchRow(Query{
		Query: "`username` = ? AND `token` = ?",
		Args:  []interface{}{u.Username, u.Token},
	})
	if login.ID <= 0 {
		return false, errors.New("user not exists")
	}

	u.ID = login.ID

	if login.ExpireAt <= time.Now().Unix() {
		return false, errors.New("token expired")
	}
	return true, nil
}

func (user *User) FetchRows(query Query) ([]User, error) {
	var list []User

	err := FetchRows(query, &list)
	if err != nil {
		return nil, errors.New("fetchrows user failed")
	}

	return list, nil
}

func (u *User) FormatJson() (UserJson, error) {
	json := UserJson{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
		AuthTime: u.AuthTime.Format("2006-01-02 15:04:05"),
	}

	return json, nil
}
