package auth

import (
	"crawler/internal/model"
	"crawler/internal/util/auth"
	"crawler/internal/util/config"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Info(c *gin.Context) {
	username, exist := c.Get("user")
	if !exist {
		req.JSON(c, req.CodeError, "sorry, fetch user failed", nil)
		return
	}

	login, err := (&model.User{}).FetchRow("`username` = ?", username.(string))
	if err != nil {
		req.JSON(c, req.CodeError, "sorry, fetch user failed", nil)
		return
	}
	js, _ := login.FormatJson()

	req.JSON(c, req.CodeSuccess, "userinfo", js)
}

func Auth(c *gin.Context) {
	cnf := config.NewConfig()
	github := &auth.GithubAuth{
		ClientId: cnf.Auth.Github.ClientId,
		ClientSecret: cnf.Auth.Github.ClientSecret,
	}

	c.Redirect(http.StatusMovedPermanently, github.RedirectAuth("http://127.0.0.1:7980/admin/callback"))
	c.Abort()
}

func Callback(c *gin.Context) {
	code := c.Query("code")

	cnf := config.NewConfig()
	github := &auth.GithubAuth{
		ClientId: cnf.Auth.Github.ClientId,
		ClientSecret: cnf.Auth.Github.ClientSecret,
	}
	token, _ := github.RequestAccessToken(code)
	gUser, _ := github.RequestUser(token)

	if gUser.Username != "aaronzjc" {
		req.JSON(c, req.CodeError, "sorry, auth failed", nil)
		return
	}

	user := model.User{
		Username: gUser.Username,
		Nickname: gUser.Username,
		Avatar: gUser.Avatar,
	}

	_ = user.Auth()

	setCookie(c, map[string]string{
		"_token": user.Token,
		"_user": user.Username,
	})

	c.Redirect(http.StatusTemporaryRedirect, "http://127.0.0.1:8080/admin")
	c.Abort()
}

func setCookie(c *gin.Context, data map[string]string) {
	cnf := config.NewConfig()
	for key, val := range data {
		c.SetCookie(key, val,  86400 * 30, "", cnf.Domain, false, false)
	}
}