package oauth

import (
	"crawler/internal/model"
	"crawler/internal/route/middleware"
	"crawler/internal/util/auth"
	"crawler/internal/util/config"
	"crawler/internal/util/req"
	"crawler/internal/util/tool"
	"fmt"
	"github.com/gin-gonic/gin"
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
	token, _ := github.RequestAccessToken(code)
	gUser, _ := github.RequestUser(token)

	if from == "admin" {
		if ok, _ := tool.ArrSearch(gUser.Username, cnf.Auth.Github.Admins); !ok {
			c.String(http.StatusForbidden, "不好意思，您没有权限。请联系管理员。")
			return
		}

		admin := model.Admin{
			Username: gUser.Username,
			Nickname: gUser.Username,
			Avatar:   gUser.Avatar,
			AuthType: model.AuthGithub,
		}
		admin.Auth()
		req.SetCookie(c, map[string]string{
			middleware.CooAdmin: admin.Username,
			middleware.CooAdminToken: admin.Token,
		})
		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", cnf.WebUrl(), "/admin"))
		c.Abort()
	} else {
		user := model.User{
			Username: gUser.Username,
			Nickname: gUser.Username,
			Avatar:   gUser.Avatar,
			AuthType: model.AuthGithub,
		}

		_ = user.Auth()

		req.SetCookie(c, map[string]string{
			middleware.CooToken: user.Token,
			middleware.CooUser: user.Username,
		})

		c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s%s", cnf.WebUrl(), "/"))
		c.Abort()
	}
}