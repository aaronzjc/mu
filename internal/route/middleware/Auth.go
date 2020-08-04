package middleware

import (
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/util/auth"
	"mu/internal/util/config"
	"mu/internal/util/logger"
	"mu/internal/util/req"
	"mu/internal/util/tool"
	"strings"
)

const (
	LoginUser = "login_user"
)

func ApiAuth(admin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if token == "" {
			req.JSON(c, req.CodeAuthFailed, "认证失败", nil)
			c.Abort()
			return
		}
		authBytes, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			req.JSON(c, req.CodeAuthFailed, "解析失败", nil)
			c.Abort()
			return
		}
		segs := strings.Split(string(authBytes), ";")
		if len(segs) != 2 || segs[0] == "" || segs[1] == "" {
			req.JSON(c, req.CodeAuthFailed, "token格式有误", nil)
			c.Abort()
			return
		}

		user := model.User{
			Username: segs[0],
			Token:    segs[1],
		}

		if _, err := user.CheckToken(); err != nil {
			logger.Info("token check failed %s .", err.Error())
			req.JSON(c, req.CodeAuthFailed, "禁止访问", nil)
			c.Abort()
			return
		}

		u, _ := user.FetchRow(model.Query{
			Query: "`id` = ?",
			Args:  []interface{}{user.ID},
		})

		if admin {
			var admins []string
			cnf := config.NewConfig()
			if u.AuthType == auth.TYPE_GITHUB {
				admins = cnf.Auth.Github.Admins
			} else if u.AuthType == auth.TYPE_WEIBO {
				admins = cnf.Auth.Weibo.Admins
			}
			if ok, _ := tool.ArrSearch(user.Username, admins); !ok {
				if ok, _ = tool.ArrSearch("everyone", admins); !ok {
					req.JSON(c, req.CodeForbidden, "不好意思，您没有权限。请联系管理员。", nil)
					c.Abort()
					return
				}
			}
		}
		c.Set(LoginUser, u.ID)

		c.Next()
	}
}
