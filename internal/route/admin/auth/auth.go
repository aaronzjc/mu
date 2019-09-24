package auth

import (
	"crawler/internal/model"
	"crawler/internal/util/auth"
	"crawler/internal/util/config"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
)

func Auth(c *gin.Context) {
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
	}

	user := model.User{
		Username: gUser.Username,
		Nickname: gUser.Username,
		Avatar: gUser.Avatar,
	}
	_ = user.Auth()

	res, _ := user.FormatJson()

	req.JSON(c, req.CodeSuccess, "成功", res)
}