package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mu/internal/util/config"
	"mu/internal/util/logger"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
)

type WeiboAccessToken struct {
	AccessToken string `json:"access_token"`
	RemindIn    string `json:"remind_in"`
	ExpiresIn   int    `json:"expires_in"`
}

type WeiboUser struct {
	ID       int64  `json:"id"`
	Username string `json:"idstr"`
	Nickname string `json:"screen_name"`
	Avatar   string `json:"avatar_large"`
}

type WeiboAuth struct {
	Auth
	ClientId     string
	ClientSecret string
}

const BY_WEIBO = "weibo"

func (auth WeiboAuth) Type() int8 {
	return TYPE_WEIBO
}

func (auth WeiboAuth) RedirectAuth() string {
	cnf := config.NewConfig()
	callback := fmt.Sprintf("%s%s", cnf.ServerUrl(), "/oauth/callback")
	url := fmt.Sprintf("https://api.weibo.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s", auth.ClientId, callback)
	return url
}

func (auth WeiboAuth) RequestAccessToken(code string) (string, error) {
	api := "https://api.weibo.com/oauth2/access_token"

	client := &http.Client{}
	resp, err := client.PostForm(api, url.Values{
		"client_id":     {auth.ClientId},
		"client_secret": {auth.ClientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {"http://api.memosa.local:7980/oauth/callback"},
	})
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

func (auth WeiboAuth) RequestUser(token string) (AuthUser, error) {
	var err error
	var api string

	// 根据access_token获取UID
	api = fmt.Sprintf("https://api.weibo.com/2/account/get_uid.json?access_token=%s", token)
	resp, err := http.Get(api)
	defer resp.Body.Close()
	if err != nil {
		logger.Error("error get weibo uid e = " + err.Error())
		return AuthUser{}, errors.New("requestUid api failed")
	}
	body, _ := ioutil.ReadAll(resp.Body)

	var uidRes map[string]interface{}
	err = json.Unmarshal(body, &uidRes)
	if err != nil {
		logger.Error("error get weibo uid e = %v", err)
		return AuthUser{}, errors.New("decode uid failed")
	}

	uid := uidRes["uid"]

	switch uid.(type) {
	case float64:
		api = fmt.Sprintf("%s?access_token=%s&uid=%s", "https://api.weibo.com/2/users/show.json", token, strconv.FormatFloat(uid.(float64), 'f', 0, 64))
	case string:
		api = fmt.Sprintf("%s?access_token=%s&uid=%s", "https://api.weibo.com/2/users/show.json", token, uid.(string))
	default:
		return AuthUser{}, errors.New("unknown type")
	}

	respUsr, err := http.Get(api)
	defer respUsr.Body.Close()
	if err != nil {
		logger.Error("request user failed %s .", err.Error())
		return AuthUser{}, errors.New("request weibo user failed")
	}
	body, err = ioutil.ReadAll(respUsr.Body)

	fmt.Println(string(body), err, api)
	var u WeiboUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logger.Error("weibo user decode failed src = %s, err = %s .", string(body), err.Error())
	}

	if u.Username == "" || u.ID == 0 {
		return AuthUser{}, errors.New("exception")
	}

	us := AuthUser{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}
	fmt.Println(u)

	return us, nil
}
