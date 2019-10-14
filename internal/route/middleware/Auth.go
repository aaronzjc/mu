package middleware

import (
	"crawler/internal/model"
	"crawler/internal/util/logger"
	"crawler/internal/util/req"
	"github.com/gin-gonic/gin"
)

const (
	CooAdmin = "_admin"
	CooAdminToken = "_admin_token"
	CooUser = "_user"
	CooToken = "_token"

	LoginUser = "login_user"
)

func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(CooAdminToken)
		username, _ := c.Cookie(CooAdmin)

		admin := &model.Admin{
			Username: username,
		}

		if pass := admin.CheckToken(token); !pass {
			logger.Info("admin token check failed .")
			req.JSON(c, req.CodeForbidden, "禁止访问", nil)
			c.Abort()
		}

		c.Next()
	}
}

func ApiAuth(forceLogin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie(CooToken)
		username, _ := c.Cookie(CooUser)

		user := model.User{
			Username: username,
			Token:    token,
		}

		if _, err := user.CheckToken(); err != nil {
			if forceLogin {
				logger.Info("token check failed %s .", err.Error())
				req.JSON(c, req.CodeForbidden, "禁止访问", nil)
				c.Abort()
			}
		} else {
			u, _ := user.FetchRow(model.Query{
				Query: "`username` = ?",
				Args: []interface{}{user.Username},
			})
			c.Set(LoginUser, u.ID)
		}

		c.Next()
	}
}
