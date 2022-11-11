package api

import (
	"net/http"
	"time"

	"github.com/aaronzjc/mu/internal/api/res"

	"github.com/gin-gonic/gin"
)

func Resp(ctx *gin.Context, code int, msg string, data interface{}) {
	if data == nil {
		data = make(map[string]struct{})
	}
	ctx.JSON(http.StatusOK, &res.RespSt{
		Code: code,
		Msg:  msg,
		Data: data,
	})
}

func SetCookies(ctx *gin.Context, data map[string]string, domain string) {
	for k, v := range data {
		ctx.SetCookie(k, v, int(time.Hour*24*30), "", domain, false, false)
	}
}
