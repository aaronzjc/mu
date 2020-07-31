package oauth

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/util/auth"
	"mu/internal/util/config"
	"mu/internal/util/req"
	"net/http"
)

func Config(c *gin.Context) {
	from := c.Query("from")
	if from != "admin" {
		from = "index"
	}

	path := "/oauth/auth?from=" + from

	cnf := config.NewConfig()
	authCnf := []map[string]string{
		{
			"name": "Github登录",
			"url":  fmt.Sprintf("%s%s", cnf.ServerUrl(), path),
		},
		{
			"name": "微博登录",
			"url": fmt.Sprintf("%s%s", cnf.ServerUrl(), path),
		},
	}

	req.JSON(c, req.CodeSuccess, "成功", authCnf)
}

func Auth(c *gin.Context) {
	from := c.Query("from")

	// 设置登录后跳转
	req.SetCookie(c, map[string]string{
		"from": from,
	})

	cnf := config.NewConfig()
	github := &auth.GithubAuth{
		ClientId:     cnf.Auth.Github.ClientId,
		ClientSecret: cnf.Auth.Github.ClientSecret,
	}

	c.Redirect(http.StatusTemporaryRedirect, github.RedirectAuth())
	c.Abort()
}

func Callback(c *gin.Context) {
	code := c.Query("code")
	from, exist := c.Cookie("from")
	if exist != nil || from != "admin" {
		// 默认跳首页
		from = "index"
	}

	if code == "" {
		c.String(http.StatusForbidden, "auth失败")
		return
	}

	cnf := config.NewConfig()
	github := &auth.GithubAuth{
		ClientId:     cnf.Auth.Github.ClientId,
		ClientSecret: cnf.Auth.Github.ClientSecret,
	}
	accessToken, _ := github.RequestAccessToken(code)
	gUser, _ := github.RequestUser(accessToken)

	if accessToken == "" || gUser.Username == "" {
		c.String(http.StatusForbidden, "获取oauth信息失败")
		return
	}

	user := model.User{
		Username: gUser.Username,
		Nickname: gUser.Username,
		Avatar:   gUser.Avatar,
		AuthType: model.AuthGithub,
	}

	_ = user.Auth()
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s;%s", user.Username, user.Token)))

	var redirect = ""
	if from == "admin" {
		redirect = fmt.Sprintf("%s?token=%s", cnf.AdminUrl(), token)
	} else {
		redirect = fmt.Sprintf("%s?token=%s", cnf.IndexUrl(), token)
	}

	c.Redirect(http.StatusTemporaryRedirect, redirect)
	c.Abort()
}
