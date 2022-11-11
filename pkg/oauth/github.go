package oauth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
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

type GithubOAuth struct {
	clientId     string
	clientSecret string
	callback     string
}

const OauthGithub = "github"

var _ OAuth = GithubOAuth{}

func NewGithubOauth(clientId string, clientSecret string, callback string) *GithubOAuth {
	return &GithubOAuth{
		clientId:     clientId,
		clientSecret: clientSecret,
		callback:     callback,
	}
}

func (auth GithubOAuth) Type() string {
	return OauthGithub
}

func (auth GithubOAuth) RedirectAuth() string {
	url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", auth.clientId, auth.callback)
	return url
}

func (auth GithubOAuth) RequestAccessToken(code string) (string, error) {
	api := "https://github.com/login/oauth/access_token"

	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&code=%s", api, auth.clientId, auth.clientSecret, code)
	client := &http.Client{
		Timeout: time.Second * 3,
	}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return "", errors.New("RequestAccessToken failed")
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var data GithubAccessToken
	err = json.Unmarshal(body, &data)
	if err != nil || data.AccessToken == "" {
		return "", errors.New("RequestAccessToken decode json failed")
	}

	return data.AccessToken, nil
}

func (auth GithubOAuth) RequestUser(token string) (User, error) {
	var err error

	url := "https://api.github.com/user"

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		return User{}, errors.New("RequestUser failed")
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var u GithubUser
	err = json.Unmarshal(body, &u)
	if err != nil {
		return User{}, err
	}

	au := User{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Avatar:   u.Avatar,
	}

	return au, nil
}
