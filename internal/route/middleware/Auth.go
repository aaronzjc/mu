package middleware

import (
	"crawler/internal/model"
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

		if _, err := user.CheckToken(); err != nil {
			log.Printf("[info] token check failed %s\n", err.Error())
			url := fmt.Sprintf("http://127.0.0.1:7980/admin/login")
			req.JSON(c, req.CodeForbidden, "禁止访问", map[string]interface{}{
				"url": url,
			})
			c.Abort()
		}

		c.Set("user", user.Username)

		c.Next()
	}
}