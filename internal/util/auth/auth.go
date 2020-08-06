package auth

import (
	"fmt"
	"mu/internal/util/config"
)

type Auth interface {
	Type() int8
	RedirectAuth() string
	RequestAccessToken(string) (string, error)
	RequestUser(string) (AuthUser, error)
}

type AuthUser struct {
	ID       int64
	Username string
	Nickname string
	Avatar   string
}

const (
	TYPE_GITHUB = iota
	TYPE_WEIBO
)

func New(by string) Auth {
	var ath Auth
	cnf := config.NewConfig()
	switch by {
	case BY_GITHUB:
		ath = GithubAuth{
			ClientId:     cnf.Auth.Github.ClientId,
			ClientSecret: cnf.Auth.Github.ClientSecret,
		}
	case BY_WEIBO:
		ath = WeiboAuth{
			ClientId:     cnf.Auth.Weibo.ClientId,
			ClientSecret: cnf.Auth.Weibo.ClientSecret,
		}
	default:
		return nil
	}

	return ath
}

func AvailableWays(from string) []map[string]string {
	cnf := config.NewConfig()
	return []map[string]string{
		{
			"name": "Github登录",
			"type": BY_GITHUB,
			"url":  fmt.Sprintf("%s?from=%s&by=%s", cnf.ServerUrl(), from, BY_GITHUB),
		},
		{
			"name": "微博登录",
			"type": BY_WEIBO,
			"url":  fmt.Sprintf("%s?from=%s&by=%s", cnf.ServerUrl(), from, BY_WEIBO),
		},
	}
}
