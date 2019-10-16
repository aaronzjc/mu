package auth

import (
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/route/middleware"
	"mu/internal/util/req"
)

func Info(c *gin.Context) {
	userId, exist := c.Get(middleware.LoginUser)
	if !exist {
		req.JSON(c, req.CodeError, "获取信息失败", nil)
		return
	}

	login, err := (&model.User{}).FetchRow(model.Query{
		Query: "`id` = ?",
		Args:  []interface{}{userId},
	})
	if err != nil {
		req.JSON(c, req.CodeError, "sorry, fetch user failed", nil)
		return
	}
	js, _ := login.FormatJson()

	req.JSON(c, req.CodeSuccess, "userinfo", js)
}

func Logout(c *gin.Context) {
	req.ClearCookie(c, []string{middleware.CooUser, middleware.CooToken})
	req.JSON(c, req.CodeSuccess, "logout success", nil)
}
