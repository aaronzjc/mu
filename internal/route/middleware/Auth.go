package middleware

import (
	"crawler/internal/util/config"
	"crawler/internal/util/req"
	"fmt"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie("_token")
		if err != nil || token == "" {
			clientId := config.NewConfig().Auth.Github.ClientId
			url := fmt.Sprintf("https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s", clientId, "http://127.0.0.1:7980/auth")
			req.JSON(c, req.CodeForbidden, "禁止访问", map[string]interface{}{
				"url": url,
			})
			c.Abort()
		}

		c.Next()
	}
}
