package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

// 给静态资源添加一个客户端缓存时间
func AddCacheControlHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		if strings.HasPrefix(c.Request.RequestURI, "/static/") {
			c.Header("Cache-Control", "max-age=31536000")
		}
	}
}
