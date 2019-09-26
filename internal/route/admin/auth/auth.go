package auth

import (
	"crawler/internal/model"
	"crawler/internal/util/auth"
	"crawler/internal/util/config"
	"crawler/internal/util/req"
	"crawler/internal/util/tool"
	"fmt"
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
		ClientId:     cnf.Auth.Github.ClientId,
		ClientSecret: cnf.Auth.Github.ClientSecret,
	}

	c.Redirect(http.StatusMovedPermanently, github.RedirectAuth())
	c.Abort()
}

func Callback(c *gin.Context) {
	code := c.Query("code")

	cnf := config.NewConfig()
	github := &auth.GithubAuth{
		ClientId:     cnf.Auth.Github.ClientId,
		ClientSecret: cnf.Auth.Github.ClientSecret,
	}
	token, _ := github.RequestAccessToken(code)
	gUser, _ := github.RequestUser(token)

	if idx := tool.ArrSearch(gUser.Username, cnf.Auth.Github.Admins); idx == -1 {
		c.String(http.StatusForbidden, "不好意思，您没有权限。请联系管理员。")
		return
	}

	user := model.User{
		Username: gUser.Username,
		Nickname: gUser.Username,
		Avatar:   gUser.Avatar,
		AuthType: model.AuthGithub,
	}

	_ = user.Auth()

	setCookie(c, map[string]string{
		"_token": user.Token,
		"_user":  user.Username,
	})

	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", cnf.WebUrl(), "/admin"))
	c.Abort()
}

func Logout(c *gin.Context) {
	clearCookie(c, []string{"_token", "_user"})
	cnf := config.NewConfig()
	c.Redirect(http.StatusMovedPermanently, fmt.Sprintf("%s%s", cnf.WebUrl(), "/admin/login"))
}

func setCookie(c *gin.Context, data map[string]string) {
	cnf := config.NewConfig()
	for key, val := range data {
		c.SetCookie(key, val, 86400*30, "", cnf.Server.Host, false, false)
	}
}

func clearCookie(c *gin.Context, keys []string) {
	cnf := config.NewConfig()
	for _, val := range keys {
		c.SetCookie(val, "", -1, "", cnf.Server.Host, false, false)
	}
}
