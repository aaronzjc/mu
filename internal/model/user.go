package model

import (
	"crawler/internal/util/tool"
	"errors"
	"fmt"
	"log"
	"time"
)

type User struct {
	ID 		int 		`gorm:"id"`
	Username string 	`gorm:"username"`
	Nickname string 	`gorm:"nickname"`
	Avatar	string 		`gorm:"avatar"`
	AuthTime string 	`gorm:"auth_time"`
	Token 	string 		`gorm:"token"`
	ExpireAt int64 		`gorm:"expire_at"`
}

type UserJson struct {
	ID 		int 		`json:"id"`
	Username string 	`json:"username"`
	Nickname string 	`json:"nickname"`
	Avatar	string 		`json:"avatar"`
	AuthTime string 	`json:"auth_time"`
	Token 	string 		`json:"token"`
	ExpireAt int64 		`json:"expire_at"`
}

func (u *User) TableName() string {
	return "user"
}

func (u *User) FetchRow(query string, args ...interface{}) (User, error) {
	var tmp User
	db := DPool().Conn
	db = db.Where(query, args...).First(&tmp)
	if err := db.Error; err != nil && !db.RecordNotFound() {
		log.Printf("[error] FetchRow err %v, exp %s \n", err, db.QueryExpr())
		return User{}, errors.New("fetch user info failed")
	}

	return tmp, nil
}

func (u *User) Create() error {
	tmp, err := u.FetchRow("`username` = ?", u.Username)
	if err != nil {
		return errors.New("create user err")
	}

	if tmp.ID > 0 {
		return errors.New(fmt.Sprintf("user with %s exists", u.Username))
	}

	db := DPool().Conn
	db = db.Create(&u)
	if err = db.Error; err != nil {
		log.Printf("[error] create err %v, exp %s \n", err, db.QueryExpr())
		return errors.New("create user err")
	}

	return nil
}

func (u *User) Update(data map[string]interface{}) error {
	db := DPool().Conn

	db = db.Model(&u).Update(data)
	if err := db.Error; err != nil {
		log.Printf("[error] update err %v, exp %s \n", err, db.QueryExpr())
		return errors.New("update user failed")
	}

	return nil
}

func (u *User) RefreshToken() error {
	token := tool.GenerateToken(u.Username)
	expireAt := time.Now().Add(time.Hour * 24 * 30).Unix()

	if err := u.Update(map[string]interface{}{
		"token": token,
		"expire_at": expireAt,
	}); err != nil {
		return errors.New("refresh token error " + err.Error())
	}

	u.Token = token
	u.ExpireAt = expireAt

	return nil
}

func (u *User) Auth() error {
	tmp, _ := u.FetchRow("`username` = ?", u.Username)
	if tmp.ID > 0 {
		err := tmp.RefreshToken()
		if err != nil {
			return errors.New("auth failed " + err.Error())
		}
		u.ID = tmp.ID
		u.AuthTime = tmp.AuthTime
		u.Token = tmp.Token
		u.ExpireAt = tmp.ExpireAt
	} else {
		u.AuthTime = time.Now().Format("2006-01-02 15:04:05")
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
	login, _ := u.FetchRow("`username` = ? AND `token` = ?", u.Username, u.Token)
	if login.ID <= 0 {
		return false, errors.New("user not exists")
	}
	if login.ExpireAt <= time.Now().Unix() {
		return false, errors.New("token expired")
	}
	return true, nil
}

func (u *User) FormatJson() (UserJson, error) {
	json := UserJson{
		ID: u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar: u.Avatar,
		AuthTime: u.AuthTime,
		Token: u.Token,
		ExpireAt: u.ExpireAt,
	}

	return json, nil
}