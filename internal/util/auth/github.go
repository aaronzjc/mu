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

type GithubAccessToken struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

type GithubUser struct {
	ID       int64  `json:"id"`
	Username string `json:"login"`
	Nickname string `json:"name"`
	Avatar   string `json:"avatar_url"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

type GithubAuth struct {
	Auth
	ClientId     string
	ClientSecret string
}

const BY_GITHUB = "github"

func (auth GithubAuth) Type() int8 {
	return TYPE_GITHUB
}

func (auth GithubAuth) RedirectAuth() string {
	cnf := config.NewConfig()
	callback := fmt.Sprintf("%s%s", cnf.ServerUrl(), "/oauth/callback")
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", auth.ClientId, callback)
	return url
}

func (auth GithubAuth) RequestAccessToken(code string) (string, error) {
	api := "https://github.com/login/oauth/access_token"

	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s", api, auth.ClientId, auth.ClientSecret, code)
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("github request access_token failed %s .", err.Error())
		return "", errors.New("RequestAccessToken failed")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	var data GithubAccessToken
	err = json.Unmarshal(body, &data)
	if err != nil {
		return "", errors.New("RequestAccessToken decode json failed")
	}

	return data.AccessToken, nil
}

func (auth GithubAuth) RequestUser(token string) (AuthUser, error) {
	var err error

	url := "https://api.github.com/user"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		logger.Error("github request user failed %s .", err.Error())
		return AuthUser{}, errors.New("RequestUser failed")
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var u GithubUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		logger.Error("github user decode failed %s .", err.Error())
	}

	au := AuthUser{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}

	return au, nil
}
