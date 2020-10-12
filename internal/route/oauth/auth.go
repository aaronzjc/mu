package oauth

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/util/auth"
	"mu/internal/util/config"
	"mu/internal/util/logger"
	"mu/internal/util/req"
	"net/http"
	"net/url"
)

func Config(c *gin.Context) {
	from := c.Query("from")
	if from != "admin" {
		from = "index"
	}

	req.JSON(c, req.CodeSuccess, "成功", auth.AvailableWays(from))
}

func Auth(c *gin.Context) {
	from := c.Query("from")
	by := c.Query("by")

	// 设置登录后跳转
	req.SetCookie(c, map[string]string{
		"from": from,
		"by":   by,
	})

	ath := auth.New(by)
	if ath == nil {
		c.Abort()
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, ath.RedirectAuth())
	c.Abort()
	return
}

func Callback(c *gin.Context) {
	var from, by string
	var err error
	code := c.Query("code")
	from, err = c.Cookie("from")
	if err != nil || from != "admin" {
		// 默认跳首页
		from = "index"
	}
	by, err = c.Cookie("by")
	if err != nil {
		c.String(http.StatusForbidden, "auth失败")
		return
	}

	if code == "" {
		c.String(http.StatusForbidden, "code缺失")
		return
	}

	ath := auth.New(by)
	if ath == nil {
		c.String(http.StatusForbidden, "未识别的回调方式")
		return
	}

	// 获取access_token
	accessToken, err := ath.RequestAccessToken(code)
	if err != nil {
		c.String(http.StatusForbidden, "获取access_token失败")
		return
	}

	// 根据access_token获取当前用户信息
	usr, err := ath.RequestUser(accessToken)
	if err != nil {
		logger.Error("request user error . e = %v", err)
		c.String(http.StatusForbidden, "获取oauth用户信息失败")
		return
	}

	user := model.User{
		Username: usr.Username,
		Nickname: usr.Nickname,
		Avatar:   usr.Avatar,
		AuthType: ath.Type(),
	}

	_ = user.Auth()
	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s;%s", user.Username, user.Token)))
	tokenEncode := url.QueryEscape(token)
	cnf := config.NewConfig()

	var redirect = ""
	if from == "admin" {
		redirect = fmt.Sprintf("%s?token=%s", cnf.AdminUrl(), tokenEncode)
	} else {
		redirect = fmt.Sprintf("%s?token=%s", cnf.IndexUrl(), tokenEncode)
	}

	c.Redirect(http.StatusTemporaryRedirect, redirect)
	c.Abort()
	return
}
