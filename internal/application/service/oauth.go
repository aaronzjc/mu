package service

import (
	"errors"
	"fmt"

	"github.com/aaronzjc/mu/internal/application/dto"
	"github.com/aaronzjc/mu/internal/config"
	"github.com/aaronzjc/mu/pkg/oauth"
)

type OAuthService interface {
	Platforms(string) []dto.OAuthPlatform
	GetPlatform(string) oauth.OAuth
	Redirect(string) string
	Auth(string, string) (oauth.User, error)
}

type OAuthServiceImpl struct {
}

var _ OAuthService = &OAuthServiceImpl{}

func (s *OAuthServiceImpl) GetPlatform(t string) oauth.OAuth {
	conf, ok := config.Get().OAuth[t]
	if !ok {
		return nil
	}
	callback := fmt.Sprintf("%s%s", config.Get().ServerUrl(), "/auth/callback")
	if t == oauth.OauthGithub {
		return oauth.NewGithubOauth(conf.ClientId, conf.ClientSecret, callback)
	}
	if t == oauth.OauthWeibo {
		return oauth.NewWeiboOauth(conf.ClientId, conf.ClientSecret, callback)
	}
	return nil
}
func (s *OAuthServiceImpl) Platforms(from string) []dto.OAuthPlatform {
	path := "/auth/redirect" // 重定向地址
	cnf := config.Get()
	return []dto.OAuthPlatform{
		{
			Name: "Github登录",
			Type: oauth.OauthGithub,
			Url:  fmt.Sprintf("%s%s?from=%s&by=%s", cnf.ServerUrl(), path, from, oauth.OauthGithub),
		},
		{
			Name: "微博登录",
			Type: oauth.OauthWeibo,
			Url:  fmt.Sprintf("%s%s?from=%s&by=%s", cnf.ServerUrl(), path, from, oauth.OauthWeibo),
		},
	}
}

func (s *OAuthServiceImpl) Redirect(t string) string {
	var a oauth.OAuth
	if a = s.GetPlatform(t); a == nil {
		return ""
	}
	return a.RedirectAuth()
}

func (s *OAuthServiceImpl) Auth(t string, code string) (oauth.User, error) {
	var a oauth.OAuth
	if a = s.GetPlatform(t); a == nil {
		return oauth.User{}, errors.New("invalid platform")
	}
	accessToken, err := a.RequestAccessToken(code)
	if err != nil {
		return oauth.User{}, errors.New("get access token err")
	}
	usr, err := a.RequestUser(accessToken)
	if err != nil {
		return oauth.User{}, errors.New("get oauth user err")
	}
	return usr, nil
}

func NewOAuthService() *OAuthServiceImpl {
	return &OAuthServiceImpl{}
}
