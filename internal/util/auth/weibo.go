package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mu/internal/util/config"
	"mu/internal/util/logger"
	"net/http"
)

type WeiboAccessToken struct {
	AccessToken string `json:"access_token"`
	RemindIn    int `json:"remind_in"`
	ExpiresIn   int `json:"expires_in"`
}

type WeiboUser struct {
	ID       int64  `json:"id"`
	Username string `json:"screen_name"`
	Avatar   string `json:"avatar_large"`
}

type WeiboAuth struct {
	ClientId     string
	ClientSecret string
}

func (auth *WeiboAuth) RedirectAuth() string {
	cnf := config.NewConfig()
	callback := fmt.Sprintf("%s%s", cnf.ServerUrl(), "/oauth/callback")
	url := fmt.Sprintf("https://api.weibo.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s", auth.ClientId, callback)
	return url
}

func (auth *WeiboAuth) RequestAccessToken(code string) (string, error) {
	api := "https://api.weibo.com/oauth2/access_token"

	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&grant_type=authorization_code&code=%s", api, auth.ClientId, auth.ClientSecret, code)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("weibo request access_token failed %s .", err.Error())
		return "", errors.New("RequestAccessToken failed")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data WeiboAccessToken
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.New("RequestAccessToken decode json failed")
	}

	return data.AccessToken, nil
}

func (auth *WeiboAuth) RequestUser(token string) (AuthUser, error) {
	var err error
	var uid string

	// 根据access_token获取UID
	resp, err := http.Get(fmt.Sprintf("https://api.weibo.com/2/account/get_uid.json?access_token=%s", token))
	defer resp.Body.Close()
	if err != nil {
		logger.Error("error get weibo uid")
		return AuthUser{}, errors.New("requestUid api failed")
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var uidRes map[string]string
	err = json.Unmarshal(body, &uidRes)
	if err != nil {
		logger.Error("error get weibo uid")
		return AuthUser{}, errors.New("decode uid failed")
	}

	uid = uidRes["uid"]

	// 根据UID获取信息
	resp, err = http.Get(fmt.Sprintf("%s?access_token=%s&uid=%s", "https://api.weibo.com/2/users/show.json", token, uid))
	if err != nil {
		logger.Error("RequestUser failed %s .", err.Error())
		return AuthUser{}, errors.New("request weibo user failed")
	}
	body, _ = ioutil.ReadAll(resp.Body)

	var u WeiboUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logger.Error("github user decode failed %s .", err.Error())
	}

	us := AuthUser{
		ID:       u.ID,
		Username: u.Username,
		Avatar:   u.Avatar,
	}

	return us, nil
}
