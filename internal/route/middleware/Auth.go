package middleware

import (
	"crawler/internal/model"
	"crawler/internal/util/config"
	"crawler/internal/util/req"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, _ := c.Cookie("_token")
		username, _ := c.Cookie("_user")

		user := model.User{
			Username: username,
			Token: token,
		}

		cnf := config.NewConfig()
		if _, err := user.CheckToken(); err != nil {
			log.Printf("[info] token check failed %s\n", err.Error())
			url := fmt.Sprintf("%s%s", cnf.ServerUrl(), "/admin/login")
			req.JSON(c, req.CodeForbidden, "禁止访问", map[string]interface{}{
				"url": url,
			})
			c.Abort()
		}

		c.Set("user", user.Username)

		c.Next()
	}
}