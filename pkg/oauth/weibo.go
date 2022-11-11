package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

type WeiboOAuth struct {
	clientId     string
	clientSecret string
	callback     string
}

const OauthWeibo = "weibo"

var _ OAuth = WeiboOAuth{}

func NewWeiboOauth(clientId string, clientSecret string, callback string) *WeiboOAuth {
	return &WeiboOAuth{
		clientId:     clientId,
		clientSecret: clientSecret,
		callback:     callback,
	}
}

func (auth WeiboOAuth) Type() string {
	return OauthWeibo
}

func (auth WeiboOAuth) RedirectAuth() string {
	url := fmt.Sprintf("https://api.weibo.com/oauth2/authorize?client_id=%s&response_type=code&redirect_uri=%s", auth.clientId, auth.callback)
	return url
}

func (auth WeiboOAuth) RequestAccessToken(code string) (string, error) {
	api := "https://api.weibo.com/oauth2/access_token"

	client := &http.Client{}
	resp, err := client.PostForm(api, url.Values{
		"client_id":     {auth.clientId},
		"client_secret": {auth.clientSecret},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {auth.callback},
	})
	if err != nil {
		return "", errors.New("RequestAccessToken failed")
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data WeiboAccessToken
	err = json.Unmarshal(body, &data)
	if err != nil || data.AccessToken == "" {
		return "", errors.New("RequestAccessToken decode json failed")
	}

	return data.AccessToken, nil
}

func (auth WeiboOAuth) RequestUser(token string) (User, error) {
	var err error
	var api string

	// 根据access_token获取UID
	api = fmt.Sprintf("https://api.weibo.com/2/account/get_uid.json?access_token=%s", token)
	resp, err := http.Get(api)
	if err != nil {
		return User{}, err
	}
	defer resp.Body.Close()
	if err != nil {
		return User{}, errors.New("requestUid api failed")
	}
	body, _ := io.ReadAll(resp.Body)

	var uidRes map[string]interface{}
	err = json.Unmarshal(body, &uidRes)
	if err != nil {
		return User{}, errors.New("decode uid failed")
	}

	uid := uidRes["uid"]

	if uid == nil {
		return User{}, errors.New("fetch uid failed")
	}

	switch uidt := uid.(type) {
	case float64:
		api = fmt.Sprintf("%s?access_token=%s&uid=%s", "https://api.weibo.com/2/users/show.json", token, strconv.FormatFloat(uidt, 'f', 0, 64))
	case string:
		api = fmt.Sprintf("%s?access_token=%s&uid=%s", "https://api.weibo.com/2/users/show.json", token, uidt)
	case int:
		api = fmt.Sprintf("%s?access_token=%s&uid=%s", "https://api.weibo.com/2/users/show.json", token, strconv.Itoa(uidt))
	default:
		return User{}, errors.New("unknown type")
	}

	respUsr, err := http.Get(api)
	if err != nil {
		return User{}, err
	}
	defer respUsr.Body.Close()
	if err != nil {
		return User{}, errors.New("request weibo user failed")
	}
	body, err = io.ReadAll(respUsr.Body)

	fmt.Println(string(body), err, api)
	var u WeiboUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		return User{}, err
	}

	if u.Username == "" || u.ID == 0 {
		return User{}, errors.New("exception")
	}

	us := User{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}

	return us, nil
}
