package auth

import (
	"github.com/gin-gonic/gin"
	"mu/internal/model"
	"mu/internal/route/middleware"
	"mu/internal/util/req"
)

func Info(c *gin.Context) {
	username, err := c.Cookie(middleware.CooUser)
	if err != nil {
		req.JSON(c, req.CodeError, "sorry, fetch user failed", nil)
		return
	}

	login, err := (&model.User{}).FetchRow(model.Query{
		Query: "`username` = ?",
		Args:  []interface{}{username},
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
