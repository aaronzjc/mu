package req

import (
	"github.com/gin-gonic/gin"
	"mu/internal/util/config"
	"net/http"
	"strings"
)

const (
	CodeSuccess   = 10000
	CodeError     = 10001
	CodeForbidden = 10002
)

func JSON(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func SetCookie(c *gin.Context, data map[string]string) {
	cnf := config.NewConfig()
	for key, val := range data {
		c.SetCookie(key, val, 86400*30, "", strings.Split(cnf.Server.Host, ":")[0], false, false)
	}
}

func ClearCookie(c *gin.Context, keys []string) {
	cnf := config.NewConfig()
	for _, val := range keys {
		c.SetCookie(val, " ", -1, "", strings.Split(cnf.Server.Host, ":")[0], false, false)
	}
}
